import { BehaviorSubject, Observable } from 'rxjs';

/**
 * Add or remove given className for specific elements monitored by observable
 */
export class ObservableCSSClassBinder {
  private _$lastState: BehaviorSubject<Set<number>> = new BehaviorSubject(
    this.initialValue,
  );

  private classAssignedElements: Set<Element> = new Set<Element>();

  constructor(
    public readonly appliedClassName: string,
    public readonly prefix: string,
    public readonly sourceObservable: Observable<Set<number>>,
    private readonly initialValue: Set<number>,
  ) {
    this.sourceObservable.subscribe((nextState) => {
      const lastState = this._$lastState.value;
      for (const key of nextState) {
        if (!lastState.has(key)) {
          this.addClassForElements(this.toClassName(key));
        }
      }
      const removedElements = [];
      for (const element of this.classAssignedElements) {
        if (!nextState.has(this.idFromElement(element))) {
          element.classList.remove(this.appliedClassName);
          removedElements.push(element);
        }
      }
      removedElements.forEach((a) => this.classAssignedElements.delete(a));
      this._$lastState.next(new Set(nextState));
    });
  }

  private toClassName(current: number): string {
    return `${this.prefix}-${current}`;
  }

  private idFromElement(elem: Element): number {
    let className = '';
    elem.classList.forEach((f) => {
      if (f.startsWith(this.prefix)) {
        className = f;
      }
    });
    return +className.substring(this.prefix.length + 1);
  }

  private addClassForElements(selectorClassName: string) {
    const elements = document.getElementsByClassName(selectorClassName);
    for (let i = 0; i < elements.length; i++) {
      const elem = elements.item(i)!;
      elem.classList.add(this.appliedClassName);
      this.classAssignedElements.add(elem);
    }
  }

  /**
   * Manually update the classes.
   * Should be called when scrolling event fired on virtual scroll viewport
   */
  public invalidate() {
    const lastState = this._$lastState.value;
    for (const key of lastState) {
      this.addClassForElements(this.toClassName(key));
    }
    const removedElements = [];
    for (const element of this.classAssignedElements) {
      if (!lastState.has(this.idFromElement(element))) {
        element.classList.remove(this.appliedClassName);
        removedElements.push(element);
      }
    }
    removedElements.forEach((a) => this.classAssignedElements.delete(a));
  }
}
