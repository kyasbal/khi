import { Injectable, Injector } from '@angular/core';
import { GoogleDriveDataLoaderService } from './google-drive';

@Injectable()
export class DataLoadSourceExtension {
  constructor(private _injector: Injector) {}

  public load(fileId: string): void {
    if (
      process.env['NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER'] &&
      process.env['NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER'] != 'false'
    ) {
      const loader = this._injector.get<GoogleDriveDataLoaderService>(
        GoogleDriveDataLoaderService,
      );
      loader.load(fileId);
    }
  }
}
