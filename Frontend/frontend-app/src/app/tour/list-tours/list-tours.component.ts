import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TourService } from '../tour.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../auth/auth.service';
import { KeypointService } from '../keypoint.service';

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
  currentUserRole?: string;
  isMyTours = false;

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private router: Router,
    private auth: AuthService,
    private keypointService: KeypointService
  ) {}

  ngOnInit() {
    this.route.paramMap.subscribe((pm) => {
      this.userId = pm.get('id') || undefined;
      this.updateIsMyTours();
      this.load();
    });

    // fetch current logged-in user info (id + role) and reload so role-based filters apply
    this.auth.whoAmI().subscribe({
      next: (user) => {
        this.currentUserId = user.id?.toString();
        this.currentUserRole = user.role;
        console.log('User role is: ', this.currentUserRole);
        this.updateIsMyTours();
        this.load(); // reload to apply role-based visibility
      },
      error: () => {
        this.currentUserId = undefined;
        this.currentUserRole = undefined;
        this.updateIsMyTours();
        this.load();
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
        // If not viewing the owner's tours, only show Published tours
        if (!this.isMyTours) {
          items = items.filter((t) => t.status === 1); // only Published for public lists
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

  viewTour(tourId: string) {
    if (!tourId) return;
    this.router.navigate(['/tours', 'view', tourId]);
  }

  confirmDelete(tourId: string) {
    if (!tourId) return;
    if (!confirm('Delete this tour? This action cannot be undone.')) return;
    this.tourService.delete(tourId).subscribe({
      next: () => this.load(),
      error: (err) => console.error('Failed to delete tour', err),
    });
  }

  canPublish(tour: any): boolean {
    return tour.status === 0 || tour.status === 2;
  }

  canArchive(tour: any): boolean {
    return tour.status === 1; // Published status
  }

  async isPublishEnabled(tourId: string): Promise<boolean> {
    try {
      const keypoints = await this.keypointService
        .getByTour(tourId)
        .toPromise();
      return (keypoints || []).length >= 2;
    } catch (err) {
      console.error('Failed to check keypoints count', err);
      return false;
    }
  }

  publishTour(tour: any) {
    this.isPublishEnabled(tour.id).then((canPublish) => {
      if (!canPublish) {
        alert('Tour must have at least 2 keypoints to be published.');
        return;
      }

      // Update status immediately for UI responsiveness
      tour.status = 1;
      tour.publishedAt = new Date().toISOString();

      const updatedTour = {
        ...tour,
        status: 1, // Published
        publishedAt: new Date().toISOString(),
      };

      this.tourService.update(tour.id, tour).subscribe({
        next: () => {
          console.log('Tour published successfully');
          // No need to reload since we already updated the UI
        },
        error: (err) => {
          console.error('Failed to publish tour', err);
          // Revert the UI change on error
          tour.status = 0;
          delete tour.publishedAt;
        },
      });
    });
  }

  archiveTour(tour: any) {
    // Update status immediately for UI responsiveness
    tour.status = 2;
    tour.archivedAt = new Date().toISOString();

    const updatedTour = {
      ...tour,
      status: 2, // Archived
      archivedAt: new Date().toISOString(),
    };

    this.tourService.update(tour.id, tour).subscribe({
      next: () => {
        console.log('Tour archived successfully');
        // No need to reload since we already updated the UI
      },
      error: (err) => {
        console.error('Failed to archive tour', err);
        // Revert the UI change on error
        tour.status = 1;
        delete tour.archivedAt;
      },
    });
  }
}
