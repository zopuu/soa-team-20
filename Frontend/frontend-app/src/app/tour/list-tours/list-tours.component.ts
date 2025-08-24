import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TourService } from '../tour.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../auth/auth.service';

@Component({
  selector: 'app-list-tours',
  templateUrl: './list-tours.component.html',
  styleUrls: ['./list-tours.component.css'],
})
export class ListToursComponent {
  @Input() userId?: string;

  tours: any[] = [];
  loading = false;
  error = '';
  currentUserId?: string;
  isMyTours = false;

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private router: Router,
    private auth: AuthService
  ) {}

  ngOnInit() {
    this.route.paramMap.subscribe((pm) => {
      this.userId = pm.get('id') || undefined;
      this.updateIsMyTours();
      this.load();
    });

    // fetch current logged-in user id (if any)
    this.auth.whoAmI().subscribe({
      next: (user) => {
        this.currentUserId = user.id?.toString();
        this.updateIsMyTours();
      },
      error: () => {
        this.currentUserId = undefined;
        this.updateIsMyTours();
      },
    });
  }

  private updateIsMyTours() {
    this.isMyTours =
      !!this.userId &&
      !!this.currentUserId &&
      this.userId === this.currentUserId;
  }

  private load(): void {
    this.loading = true;
    this.error = '';
    const req$ = this.userId
      ? this.tourService.getAllByUser(this.userId)
      : this.tourService.getAllTours();

    req$.subscribe({
      next: (data) => {
        let items = (data || []) as any[];
        // backend uses numeric status enum: Draft=0, Published=1, Archived=2
        if (!this.userId) {
          items = items.filter((t) => t.status === 1); // only Published on all-tours
        }
        this.tours = items;
        this.loading = false;
      },
      error: (err) => {
        console.error('Failed to load tours', err);
        this.error = 'Failed to load tours.';
        this.loading = false;
      },
    });
  }

  trackById = (_: number, t: any) => t.id;

  difficultyLabel(d: number): string {
    // backend: 0=Beginner,1=Intermediate,2=Advanced,3=Pro
    switch (d) {
      case 0:
        return 'Beginner';
      case 1:
        return 'Intermediate';
      case 2:
        return 'Advanced';
      case 3:
        return 'Pro';
      default:
        return 'Unknown';
    }
  }

  statusLabel(s: number): string {
    switch (s) {
      case 0:
        return 'Draft';
      case 1:
        return 'Published';
      case 2:
        return 'Archived';
      default:
        return 'Unknown';
    }
  }

  goToCreateKeypoint(tourId: string) {
    // navigate to route that allows creating keypoint for a tour
    this.router.navigate(['/tours', tourId, 'create-keypoint']);
  }

  confirmDelete(tourId: string) {
    if (!tourId) return;
    if (!confirm('Delete this tour? This action cannot be undone.')) return;
    this.tourService.delete(tourId).subscribe({
      next: () => this.load(),
      error: (err) => console.error('Failed to delete tour', err),
    });
  }
}
