import { Injectable } from '@angular/core';
import { combineLatest, filter, from, Subject } from 'rxjs';

/**
 * The information to be shown as a notification
 */
export interface KHINotification {
  /**
   * Title of the notification
   */
  title: string;
  /**
   * Body of the notification
   */
  body: string;
}

/**
 * NotificationManager provides features to show notifications with WebNotification API.
 * The initialize method must be called at the beginning of this app.
 */
@Injectable({
  providedIn: 'root',
})
export class NotificationManager {
  private notificationQueue: Subject<KHINotification> =
    new Subject<KHINotification>();

  /**
   * Notify the given information as the Notification.
   * It will be shown after user accepted the notification permission.
   * @param notification
   */
  notify(notification: KHINotification) {
    this.notificationQueue.next(notification);
  }

  /**
   * Initialize NotificationManager. This will request the notification permission from user.
   */
  initialize(): void {
    combineLatest([
      from(Notification.requestPermission()).pipe(
        filter((result) => result === 'granted'),
      ),
      this.notificationQueue,
    ]).subscribe(([, queue]) => {
      new Notification(queue.title, {
        body: queue.body,
        icon: 'assets/icons/khi.png',
      });
    });
  }
}
