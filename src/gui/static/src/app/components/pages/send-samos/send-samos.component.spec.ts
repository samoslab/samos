import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendSamosComponent } from './send-samos.component';

describe('SendSamosComponent', () => {
  let component: SendSamosComponent;
  let fixture: ComponentFixture<SendSamosComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendSamosComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendSamosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
