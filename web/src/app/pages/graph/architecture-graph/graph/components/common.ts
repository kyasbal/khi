import {
  GraphResourceData,
  NamespacedArchGraphResource,
} from '../../../../../common/schema/graph-schema';
import {
  AnchorPoints,
  Direction,
  ElementStyle,
  GraphObject,
} from '../base/base-containers';
import {
  $alignedGroup,
  $boxed_label,
  $empty,
  $label,
  $rect,
} from '../builder-alias';
import {
  GRAPH_COLORS,
  NAMESPACE_LABEL_BACKGROUND,
  NODE_KIND_LABEL,
} from '../styles';

export function generateKindLabel(
  label: string,
  accentColor: string,
  labelStyle: ElementStyle,
  round = 4,
): GraphObject {
  return $rect()
    .withStyle({
      fill: accentColor,
      rx: round,
      ry: round,
    })
    .withChildren([$label(label).withStyle(labelStyle).withMargin(3, 5, 3, 5)]);
}

export function generateDeletionOrUpdateLabel(
  resourceData: GraphResourceData,
  resourceColor: string,
  fontSize: number,
): GraphObject[] {
  if (!resourceData.deletedAt && !resourceData.updatedAt) return [];
  if (resourceData.updatedAt) {
    return [
      $boxed_label(
        `Updated ${resourceData.updatedAt}`,
        { fill: resourceColor, rx: 4, ry: 4 },
        { fill: 'white', 'font-size': fontSize },
        [3, 20, 3, 20],
      )
        .withAnchor(AnchorPoints.BOTTOM_RIGHT)
        .withPivot(AnchorPoints.BOTTOM_RIGHT),
    ];
  }
  return [
    $boxed_label(
      `Deleted ${resourceData.deletedAt}`,
      { fill: GRAPH_COLORS.ERROR, rx: 4, ry: 4 },
      { fill: 'white', 'font-size': fontSize },
      [3, 20, 3, 20],
    )
      .withAnchor(AnchorPoints.BOTTOM_RIGHT)
      .withPivot(AnchorPoints.BOTTOM_RIGHT),
  ];
}

export function generateDeletedResourceStyle(
  resourceData: GraphResourceData,
): ElementStyle {
  if (!resourceData.deletedAt) return {};
  return {
    'stroke-dasharray': '10 10',
  };
}

export function generateNameAndNamespaceBox(
  nameAndNamespaced: NamespacedArchGraphResource,
  size: number,
): GraphObject {
  return $alignedGroup(Direction.Vertical)
    .withGap(5)
    .withChildren([
      $empty().withChildren([
        $label(nameAndNamespaced.name)
          .withStyle({
            'font-size': size,
          })
          .withAnchor(AnchorPoints.CENTER)
          .withPivot(AnchorPoints.CENTER),
      ]),
      $empty().withChildren([
        $boxed_label(
          nameAndNamespaced.namespace,
          NAMESPACE_LABEL_BACKGROUND,
          {
            ...NODE_KIND_LABEL,
            'font-size': size * 0.8,
          },
          [3, 5, 3, 5],
        )
          .withAnchor(AnchorPoints.CENTER)
          .withPivot(AnchorPoints.CENTER),
      ]),
    ]);
}
