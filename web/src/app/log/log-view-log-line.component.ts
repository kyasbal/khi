import { Component, EventEmitter, Input, Output } from '@angular/core';
import { LogEntry } from '../store/log';

/**
 * A line of log in log view list.
 */
@Component({
  selector: 'khi-log-view-log-line',
  templateUrl: './log-view-log-line.component.html',
  styleUrls: ['./log-view-log-line.component.sass'],
})
export class LogViewLogLineComponent {
  /**
   * The LogEntry to show in this line.
   */
  @Input()
  log!: LogEntry;

  /**
   * An event triggered when user's mouse curosr hover on this line.
   */
  @Output()
  lineHover: EventEmitter<LogEntry> = new EventEmitter();

  /**
   * An event triggered when user clicked this log line.
   */
  @Output()
  lineClick: EventEmitter<LogEntry> = new EventEmitter();

  onClick() {
    this.lineClick.emit(this.log);
  }

  onHover() {
    this.lineHover.emit(this.log);
  }
}
