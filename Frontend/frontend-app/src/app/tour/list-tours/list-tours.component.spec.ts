import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListToursComponent } from './list-tours.component';

describe('ListToursComponent', () => {
  let component: ListToursComponent;
  let fixture: ComponentFixture<ListToursComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ListToursComponent]
    });
    fixture = TestBed.createComponent(ListToursComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
