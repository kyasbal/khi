#version 300 es
precision highp float;
// Fragment shader for drawing events

in vec4 eventColor;
in vec4 borderColor;
in vec4 severityColor;
in vec4 severityBorderColor;

in vec2 originalPosition;
out vec4 outColor;
flat in int selectionStatus;

const float edgeThickness = 0.35f;
const float highlightedEdgeThickness = 0.4f;
const float selectedEdgeThickness = 0.5f;

const float severityBorderHeight = 0.2f;
const float borderThicknessBetweenSeverityAndLogType = 0.4f;

void main() {
    float alpha = 0.8f;
    float bodySize = 1.0f - edgeThickness;
    if(selectionStatus == 1) {
        alpha = 0.9f;
        bodySize = 1.f - highlightedEdgeThickness;
    } else if(selectionStatus == 2) {
        alpha = 1.f;
        bodySize = 1.f - selectedEdgeThickness;
    }
    float severityBorder = severityBorderHeight - originalPosition.x;
    float severityToEventRatio = step(originalPosition.y, severityBorder);
    float borderToBodyRatio = step(max(abs(originalPosition.x), abs(originalPosition.y)), bodySize);
    outColor = mix(mix(severityBorderColor, severityColor, borderToBodyRatio), mix(borderColor, eventColor, borderToBodyRatio), severityToEventRatio);
}
