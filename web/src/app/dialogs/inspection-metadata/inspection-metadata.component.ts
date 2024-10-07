import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { InspectionMetadataOfRunResult } from '../../common/schema/api-types';

@Component({
  templateUrl: './inspection-metadata.component.html',
  styleUrls: ['./inspection-metadata.component.sass'],
})
export class InspectionMetadataDialogComponent {
  constructor(
    @Inject(MAT_DIALOG_DATA) public data: InspectionMetadataOfRunResult,
  ) {}
}
