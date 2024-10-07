import { Title } from '@angular/platform-browser';
import { RouterStateSnapshot, TitleStrategy } from '@angular/router';
import { Injectable } from '@angular/core';

/**
 * Control the window title regarding routing paths
 */
@Injectable()
export class KHITitleStrategy extends TitleStrategy {
  constructor(private readonly title: Title) {
    super();
  }

  override updateTitle(snapshot: RouterStateSnapshot): void {
    const sessionId = snapshot.root.firstChild?.params['sessionId'];
    const baseTitle = this.buildTitle(snapshot) ?? '';
    if (typeof sessionId === 'string') {
      this.title.setTitle(`${baseTitle} (${sessionId})`);
    } else {
      this.title.setTitle(baseTitle);
    }
  }
}
