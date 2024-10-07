import { CommonModule } from '@angular/common';
import { Component, Inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatInputModule } from '@angular/material/input';
import {
  debounceTime,
  map,
  ReplaySubject,
  shareReplay,
  startWith,
  switchMap,
  take,
} from 'rxjs';
import { PopupFormRequestWithClient } from 'src/app/services/popup/popup-manager';

export interface AdditionalInputPopupDialogRequest {
  formRequest: PopupFormRequestWithClient;
}

@Component({
  selector: 'khi-additional-input-popup',
  standalone: true,
  templateUrl: './additional-input-popup.component.html',
  styleUrls: ['./additional-input-popup.component.sass'],
  imports: [CommonModule, MatButtonModule, MatInputModule],
})
export class AdditionalInputPopupComponent {
  readonly formRequest: PopupFormRequestWithClient;

  validationRequests: ReplaySubject<string> = new ReplaySubject<string>(1);

  validationError = this.validationRequests.pipe(
    debounceTime(500),
    switchMap((value) =>
      this.data.formRequest.client.validate({
        id: this.data.formRequest.id,
        value,
      }),
    ),
    map((result) => result.validationError),
    shareReplay({
      bufferSize: 1,
      refCount: true,
    }),
    startWith('please input the value'),
  );

  isValid = this.validationError.pipe(map((e) => e === ''));

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: AdditionalInputPopupDialogRequest,
    private readonly dialogRef: MatDialogRef<object, void>,
  ) {
    dialogRef.disableClose = true;
    this.formRequest = data.formRequest;
  }

  onTextAreaUpdate(event: Event) {
    const textarea = event.target as HTMLTextAreaElement;
    const inputValue = textarea.value;
    this.onUpdateInput(inputValue);
  }

  onUpdateInput(inputValue: string) {
    this.validationRequests.next(inputValue);
  }

  onSubmit() {
    this.validationRequests
      .pipe(
        take(1),
        switchMap((value) =>
          this.data.formRequest.client.answer({
            id: this.data.formRequest.id,
            value,
          }),
        ),
      )
      .subscribe(() => {
        this.dialogRef.close();
      });
  }
}
