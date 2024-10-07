import { Component } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { StartupDialogComponent } from '../dialogs/startup/startup.component';

@Component({
  selector: 'khi-main-menu',
  templateUrl: './main-menu.component.html',
  styleUrls: ['./main-menu.component.sass'],
})
export class MainMenuComponent {
  constructor(
    private readonly dialog: MatDialog,
    public readonly inspectionDataStore: InspectionDataStoreService,
  ) {}

  openStartupMenu() {
    this.dialog.open(StartupDialogComponent, {
      maxWidth: '100vw',
      panelClass: 'startup-modalbox',
    });
  }
}
