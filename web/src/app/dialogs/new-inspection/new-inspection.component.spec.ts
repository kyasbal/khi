import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatOptionModule } from '@angular/material/core';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { NewInspectionDialogComponent } from './new-inspection.component';
import { MatStepperModule } from '@angular/material/stepper';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { BACKEND_API } from 'src/app/services/api/backend-api-interface';
import { BACKEND_CONNECTION } from 'src/app/services/api/backend-connection.service';
import { BackendConnectionService } from 'src/app/services/api/backend-connection-interface';

describe('NewInspectionDialogTest', () => {
  let component: NewInspectionDialogComponent;
  let fixture: ComponentFixture<NewInspectionDialogComponent>;
  let backendConnectionSpy: jasmine.SpyObj<BackendConnectionService>;

  beforeEach(async () => {
    backendConnectionSpy = jasmine.createSpyObj<BackendConnectionService>(
      'BackendConnectionService',
      ['inspectionTypes'],
    );
    await TestBed.configureTestingModule({
      declarations: [NewInspectionDialogComponent],
      imports: [
        NoopAnimationsModule,
        MatDialogModule,
        MatButtonModule,
        MatSelectModule,
        MatOptionModule,
        MatInputModule,
        MatCheckboxModule,
        ReactiveFormsModule,
        FormsModule,
        MatAutocompleteModule,
        MatStepperModule,
        MatProgressBarModule,
      ],
      providers: [
        {
          provide: MatDialogRef,
          useValue: null,
        },
        {
          provide: BACKEND_API,
          useValue: {},
        },
        {
          provide: BACKEND_CONNECTION,
          useValue: backendConnectionSpy,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(NewInspectionDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
