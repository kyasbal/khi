import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MetaTableRowComponent } from './meta-table-row.component';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';

describe('MetaTableRowComponent', () => {
  let fixture: ComponentFixture<MetaTableRowComponent>;
  let component: MetaTableRowComponent;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [MetaTableRowComponent],
      imports: [MatIconModule, MatTooltipModule, ClipboardModule],
    });

    fixture = TestBed.createComponent(MetaTableRowComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show the key given from Input', () => {
    component.key = 'key';
    fixture.detectChanges();
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.querySelector('.key')?.textContent).toContain('key');
  });

  it('should show the value given from Input', () => {
    component.value = 'value';
    fixture.detectChanges();
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.querySelector('.value-inner')?.textContent).toContain(
      'value',
    );
  });

  it('should show icon only when icon is given from Input', () => {
    component.icon = '';
    const compiled = fixture.nativeElement as HTMLElement;
    fixture.detectChanges();
    expect(compiled.querySelector('.icon')).toBeNull();

    component.icon = 'home';
    fixture.detectChanges();
    expect(compiled.querySelector('.icon')).not.toBeNull();
    expect(compiled.querySelector('.icon')?.textContent).toContain('home');
  });
});
