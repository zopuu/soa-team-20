import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PositionSimComponent } from './position-sim.component';

describe('PositionSimComponent', () => {
  let component: PositionSimComponent;
  let fixture: ComponentFixture<PositionSimComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [PositionSimComponent]
    });
    fixture = TestBed.createComponent(PositionSimComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
