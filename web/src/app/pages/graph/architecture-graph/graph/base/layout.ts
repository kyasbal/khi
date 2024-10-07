import {
  GraphObjectWithElementType,
  GraphObject,
  Size,
  ZeroSize2,
} from './base-containers';
/**
 * [<-1fr->[left-content]<-1fr->[right-content]]
 */
export class ContentMetadataPairLayout extends GraphObjectWithElementType<SVGGElement> {
  private _currentSize: Size = ZeroSize2();

  private _minGap = 0;

  constructor(
    public readonly left: GraphObject,
    public readonly right: GraphObject,
  ) {
    super(document.createElementNS('http://www.w3.org/2000/svg', 'g'));

    this.left.transform.parentGraphObject = this;
    this.right.transform.parentGraphObject = this;
    this.transform.onSvgElementUpdate = () => {
      return;
    };
    this.transform.onCustomContentSizeCalculate = () => this._currentSize;

    this.transform.onOverrideChildrenOffsets = (originalPos) => {
      const minWidth = this.transform.parent!.calculateSize().width;
      const leftSize = this.left.transform.calculateSize();
      const rightSize = this.right.transform.calculateSize();
      const contentWidth = leftSize.width + rightSize.width + this._minGap * 2;
      const contentHeight = Math.max(leftSize.height, rightSize.height);
      const containerWidth = Math.max(minWidth, contentWidth);
      this._currentSize = {
        width: containerWidth,
        height: contentHeight,
      };
      const spaceWidth = containerWidth - contentWidth;
      return [
        {
          x: originalPos.x + spaceWidth / 2,
          y: originalPos.y,
        },
        {
          x: originalPos.x + containerWidth - rightSize.width,
          y: originalPos.y,
        },
      ];
    };
  }

  public withMinimumGap(minGap: number): this {
    this._minGap = minGap;
    return this;
  }
}
