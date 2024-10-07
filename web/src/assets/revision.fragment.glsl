#version 300 es
precision highp float;
precision highp int;

layout(std140) uniform ViewState {
    // Resolution of the canvas (not viewport)
    vec2 resolution;
    // How many pixels are used for 1ms distance
    float pixelPerTime;

    float pixelScale;
    // Offset unixtime to the left most edge.
    float offsetToLeft;

    float logTypeCount;
} vs; // UBO index is 0 for vs

out vec4 outColor;

in vec2 originalPosition;
// Revision rectangle size in pixels
in vec2 actualSize;
in vec3 revisionBaseColor;
flat in int revisionIndex;
flat in int selectionStatus;

const vec2 edgeThickness = vec2(2.f, 1.f);
// To avoid using large empty space around digit in uv space, the uv will be scaled to use only the center part of this width.
const float fontExtractWidth = 0.58f;
const vec2 fontSize = vec2(fontExtractWidth, 1.0f) * .6f;
// Offset of revision indexes from bottom left
const vec2 fontPadding = vec2(12, 8);

const float loge10 = 2.302585092994046f;
const float epsilon = 0.00001f;

uniform sampler2D numberTexture;

float number(vec2 uv, int num) {
    uv.x = uv.x * fontExtractWidth + (1.f - fontExtractWidth) / 2.f;
    vec2 offset = vec2(0.1f * float(num), 0);
    vec2 size = vec2(0.1f, 1.0f);
    vec2 uvFlipY = offset + size * uv;
    return texture(numberTexture, vec2(uvFlipY.x, 1.f - uvFlipY.y)).a;
}

float log10(float x) {
    return log(x) / loge10;
}

void main() {
    // Draw border of revision rectangle
    // edgeSize is the thickness in uv coordinate.
    vec2 edgeSize = vec2(2.f / actualSize) * edgeThickness;
    // Border become 1 on fragments on the border.
    float border = max(max(step(originalPosition.x, -1.0f + edgeSize.x), step(1.f - edgeSize.x, originalPosition.x)), // horizontal edge
    max(step(originalPosition.y, -1.0f + edgeSize.y), step(1.f - edgeSize.y, originalPosition.y)));

    vec3 baseColor = revisionBaseColor;

    vec2 uvPadding = fontPadding / actualSize;
    float digitCount = floor(log10(float(max(1, revisionIndex)) + 0.001f) + 1.f); // max is needed to avoid log10(0) (NaN)
    vec2 uvSubtractingPads = (originalPosition - uvPadding) / 2.0f + 0.5f;
    vec2 uvScale = actualSize / (fontSize * actualSize.y);
    vec2 numberUv = max(uvSubtractingPads * uvScale, vec2(0));// Shrink uv to 0 for the points exceeding digit count of the uv
    float originalUvx = numberUv.x;
    int divisor = int(pow(10.f, float(int(digitCount) - int(floor(numberUv.x)) - 1)) + epsilon);
    numberUv.x = fract(numberUv.x);
    int currentDigit = revisionIndex / divisor;
    float isDigit = number(clamp(numberUv, vec2(0), vec2(1)), currentDigit % 10);
    isDigit *= step(originalUvx, digitCount);

    float baseAlpha = 0.6f;
    vec3 digitColor = vec3(0);
    if(selectionStatus == 2) {
        baseAlpha = 0.9f;
        digitColor = vec3(1);
    } else if(selectionStatus == 1) {
        baseAlpha = 0.5f;
    }

    outColor = mix(vec4(baseColor, baseAlpha + 0.2f * border), vec4(digitColor, 1), isDigit);
}
