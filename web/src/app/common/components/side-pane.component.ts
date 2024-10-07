import { Component, Input, OnChanges } from '@angular/core';
import { ResizingCalculator } from '../resizable-pane/resizing-calculator';
import { Observable } from 'rxjs';

@Component({
  selector: 'khi-side-pane',
  templateUrl: './side-pane.component.html',
  styleUrls: ['./side-pane.component.sass'],
})
export class SidePaneComponent implements OnChanges {
  static readonly DEFAULT_PANE_WIDTH = 300;

  static readonly MINIMUM_PANE_WIDTH = 100;

  @Input()
  paneTitle = '';

  @Input()
  icon = '';

  @Input()
  resizeCalculator!: ResizingCalculator;

  @Input()
  areaNameInResizer: string = '';

  areaSize!: Observable<number>;

  ngOnChanges(): void {
    this.areaSize = this.resizeCalculator.areaSize(this.areaNameInResizer);
  }

  resizeStart() {
    const resizeMove = (e: MouseEvent) => {
      const size = this.resizeCalculator.getAreaSize(this.areaNameInResizer);
      this.resizeCalculator.resizeArea(
        this.areaNameInResizer,
        size - e.movementX,
      );
    };
    window.addEventListener('mouseup', () => {
      window.removeEventListener('mousemove', resizeMove);
    });
    window.addEventListener('mousemove', resizeMove);
  }
}
