import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ForbbidenComponent } from './forbbiden.component';

describe('ForbbidenComponent', () => {
  let component: ForbbidenComponent;
  let fixture: ComponentFixture<ForbbidenComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ForbbidenComponent]
    });
    fixture = TestBed.createComponent(ForbbidenComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
