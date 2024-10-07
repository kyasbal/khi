import {
  AfterViewInit,
  Directive,
  EventEmitter,
  HostListener,
  Output,
} from '@angular/core';
import { Subject, distinctUntilChanged } from 'rxjs';

/**
 *
 */
@Directive({
  selector: '[khiCaptureShiftKey]',
})
export class CaptureShiftKeyDirective implements AfterViewInit {
  constructor() {
    this.shiftStatus
      .pipe(distinctUntilChanged())
      .subscribe((status) => this.shiftStatusChange.emit(status));
  }

  @Output() shiftStatusChange = new EventEmitter<boolean>();

  private shiftStatus = new Subject<boolean>();

  private containsMouse = false;

  @HostListener('mouseenter')
  mouseEnter() {
    this.containsMouse = true;
  }

  @HostListener('mouseleave')
  mouseLeave() {
    this.containsMouse = false;
  }

  ngAfterViewInit(): void {
    // the target element may be not able to have focus. In the case, the key related event won't be fired.
    // it will be handled when the event was propagated to the window
    window.addEventListener('keydown', (m) => {
      if (this.containsMouse) this.keyboardEvent(m);
    });
    // to support release shift key not on the target element, it needs to be handled on capture phase.
    // (For in case some element stopping the propagation)
    window.addEventListener('keyup', (m) => this.keyboardEvent(m), {
      capture: true,
    });
  }

  keyboardEvent(m: KeyboardEvent) {
    this.shiftStatus.next(m.shiftKey);
  }
}
