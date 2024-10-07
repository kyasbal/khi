import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatDialogModule, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { InspectionMetadataDialogComponent } from './inspection-metadata.component';
import { provideHttpClient } from '@angular/common/http';

describe('TaskMetadataViewDialogComponent', () => {
  let component: InspectionMetadataDialogComponent;
  let fixture: ComponentFixture<InspectionMetadataDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [InspectionMetadataDialogComponent],
      imports: [MatDialogModule],
      providers: [
        {
          provide: MAT_DIALOG_DATA,
          useValue: [],
        },
        provideHttpClient(),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(InspectionMetadataDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
