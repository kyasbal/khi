import {
  ComponentFixture,
  fakeAsync,
  TestBed,
  tick,
} from '@angular/core/testing';
import {
  AdditionalInputPopupComponent,
  AdditionalInputPopupDialogRequest,
} from './additional-input-popup.component';
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { MatDialogHarness } from '@angular/material/dialog/testing';
import { HarnessLoader } from '@angular/cdk/testing';
import { TestbedHarnessEnvironment } from '@angular/cdk/testing/testbed';
import { Component } from '@angular/core';
import { PopupFormRequestWithClient } from 'src/app/services/popup/popup-manager';
import { By } from '@angular/platform-browser';
import { MockPopupClient } from 'src/app/services/popup/mock';

describe('AdditionalInputPopupComponent in dialog context', () => {
  @Component({
    template: '<div></div>',
    standalone: true,
  })
  class TestingDialogWrapComponent {}

  let testingWrapper: ComponentFixture<TestingDialogWrapComponent>;
  let loader: HarnessLoader;
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestingDialogWrapComponent, MatDialogModule],
    }).compileComponents();
    testingWrapper = TestBed.createComponent(TestingDialogWrapComponent);
    testingWrapper.detectChanges();
    loader = await TestbedHarnessEnvironment.documentRootLoader(testingWrapper);
  });

  async function testIfDialogShowingUpWithParam(
    request: PopupFormRequestWithClient,
  ) {
    const matDialog = TestBed.inject(MatDialog);
    matDialog.open<
      AdditionalInputPopupComponent,
      AdditionalInputPopupDialogRequest
    >(AdditionalInputPopupComponent, {
      data: {
        formRequest: request,
      },
    });
    const dialogs = await loader.getAllHarnesses(MatDialogHarness);
    expect(dialogs.length).toBe(1);
  }

  it('should be instanciated with type=text', async () => {
    await testIfDialogShowingUpWithParam({
      id: 'foo',
      type: 'text',
      title: 'foo title',
      description: 'test description',
      placeholder: 'test placeholder',
      client: new MockPopupClient(),
    });
  });
});

describe('AdditionalInputPopupComponent', () => {
  let matDialogRefSpy: jasmine.SpyObj<
    MatDialogRef<AdditionalInputPopupDialogRequest, void>
  >;
  beforeEach(async () => {
    matDialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['close'], {
      disableClose: false,
    });
  });
  it('should have disbaled submit button at first', async () => {
    await TestBed.configureTestingModule({
      imports: [AdditionalInputPopupComponent, MatDialogModule],
      providers: [
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            formRequest: {
              id: 'foo',
              type: 'text',
              title: 'foo title',
              description: 'test description',
              placeholder: 'test placeholder',
              client: new MockPopupClient(),
            },
          },
        },
        {
          provide: MatDialogRef,
          useValue: matDialogRefSpy,
        },
      ],
    }).compileComponents();
    const fixture = TestBed.createComponent(AdditionalInputPopupComponent);
    fixture.detectChanges();
    const button = fixture.debugElement.query(By.css('.submit-button'));
    expect(button.nativeElement.disabled).toBe(true);
  });

  it('should update the disabled status of submit button by input', fakeAsync(async () => {
    await TestBed.configureTestingModule({
      imports: [AdditionalInputPopupComponent, MatDialogModule],
      providers: [
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            formRequest: {
              id: 'foo',
              type: 'text',
              title: 'foo title',
              description: 'test description',
              placeholder: 'test placeholder',
              client: new MockPopupClient(),
            },
          },
        },
        {
          provide: MatDialogRef,
          useValue: matDialogRefSpy,
        },
      ],
    }).compileComponents();
    const fixture = TestBed.createComponent(AdditionalInputPopupComponent);
    fixture.detectChanges();
    const textarea = fixture.debugElement.query(
      By.css('.input-text-type-textarea'),
    );
    const button = fixture.debugElement.query(By.css('.submit-button'));

    textarea.nativeElement.value = 'valid';
    textarea.nativeElement.dispatchEvent(new Event('input'));
    tick(1000);
    fixture.detectChanges();
    expect(button.nativeElement.disabled).toBe(false);

    textarea.nativeElement.value = 'invalid';
    textarea.nativeElement.dispatchEvent(new Event('input'));
    tick(1000);
    fixture.detectChanges();
    expect(button.nativeElement.disabled).toBe(true);
  }));

  it('should close dialog after submit', fakeAsync(async () => {
    await TestBed.configureTestingModule({
      imports: [AdditionalInputPopupComponent, MatDialogModule],
      providers: [
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            formRequest: {
              id: 'foo',
              type: 'text',
              title: 'foo title',
              description: 'test description',
              placeholder: 'test placeholder',
              client: new MockPopupClient(),
            },
          },
        },
        {
          provide: MatDialogRef,
          useValue: matDialogRefSpy,
        },
      ],
    }).compileComponents();
    const fixture = TestBed.createComponent(AdditionalInputPopupComponent);
    fixture.detectChanges();

    const textarea = fixture.debugElement.query(
      By.css('.input-text-type-textarea'),
    );
    const button = fixture.debugElement.query(By.css('.submit-button'));

    textarea.nativeElement.value = 'valid';
    textarea.nativeElement.dispatchEvent(new Event('input'));
    tick(1000);
    fixture.detectChanges();
    button.nativeElement.click();
    tick(1000);
    expect(matDialogRefSpy.close).toHaveBeenCalled();
  }));

  it('should show the valdiation error', fakeAsync(async () => {
    await TestBed.configureTestingModule({
      imports: [AdditionalInputPopupComponent, MatDialogModule],
      providers: [
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            formRequest: {
              id: 'foo',
              type: 'text',
              title: 'foo title',
              description: 'test description',
              placeholder: 'test placeholder',
              client: new MockPopupClient(),
            },
          },
        },
        {
          provide: MatDialogRef,
          useValue: matDialogRefSpy,
        },
      ],
    }).compileComponents();
    const fixture = TestBed.createComponent(AdditionalInputPopupComponent);
    fixture.detectChanges();

    const textarea = fixture.debugElement.query(
      By.css('.input-text-type-textarea'),
    );
    const validationError = fixture.debugElement.query(
      By.css('.validation-error'),
    );

    textarea.nativeElement.value = 'invalid';
    textarea.nativeElement.dispatchEvent(new Event('input'));
    tick(1000);
    fixture.detectChanges();
    expect(validationError.nativeElement.textContent).toBe(
      "invalid isn't valid",
    );
  }));
});
