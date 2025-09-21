import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TourService } from '../tour.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../auth/auth.service';
import { KeypointService } from '../keypoint.service';
import { Tour, TourStatus, TransportType } from '../tour.model';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

interface TourRating {
  id?: string;
  tourId?: string;
  rating: number;
  comment?: string;
  touristName?: string;
  touristEmail?: string;
  visitedAt?: string | Date;
  commentedAt?: string | Date;
  createdAt?: string | Date;
  images?: string[];
}

@Component({
  selector: 'app-list-tours',
  templateUrl: './list-tours.component.html',
  styleUrls: ['./list-tours.component.css'],
})
export class ListToursComponent {
  showRatingsModal = false;
  ratingsLoading = false;
  ratingsError = '';
  ratings: TourRating[] = [];
  selectedTourForRatings: string | null = null;

  // Helper functions to convert string enums to numbers for backend
  private statusToNumber(status: string): number {
    switch (status) {
      case 'Draft':
        return 0;
      case 'Published':
        return 1;
      case 'Archived':
        return 2;
      default:
        return 0;
    }
  }

  private difficultyToNumber(difficulty: string): number {
    switch (difficulty) {
      case 'Beginner':
        return 0;
      case 'Intermediate':
        return 1;
      case 'Advanced':
        return 2;
      case 'Pro':
        return 3;
      default:
        return 0;
    }
  }

  private transportTypeToNumber(transportType: string): number {
    switch (transportType) {
      case 'Walking':
        return 0;
      case 'Bicycle':
        return 1;
      case 'Bus':
        return 2;
      default:
        return 0;
    }
  }
  @Input() userId?: string;

  tours: any[] = [];
  loading = false;
  error = '';
  currentUserId?: string;
  currentUserRole?: string;
  isMyTours = false;

  showRateModal = false;
  selectedTourId?: string;
  reviewForm!: FormGroup;
  selectedImages: File[] = [];
  submittingReview = false;

  // (opciono) auto-popuna iz whoAmI
  userFullName?: string;
  userEmail?: string;


  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private router: Router,
    private auth: AuthService,
    private keypointService: KeypointService,
    private fb: FormBuilder,
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

rateTour(tourId: string) {
  if (!tourId) return;
  this.selectedTourId = tourId;

  const todayISO = new Date().toISOString().slice(0, 10); // yyyy-mm-dd

  this.reviewForm = this.fb.group({
    rating: [5],
    comment: [''],
    touristName: [this.userFullName ?? ''],
    touristEmail: [this.userEmail ?? ''],
    visitedAt: [todayISO],
    commentedAt: [todayISO],

    imageUrl1: ['', [Validators.pattern('https?://.+')]],
    imageUrl2: ['', [Validators.pattern('https?://.+')]],
    imageUrl3: ['', [Validators.pattern('https?://.+')]],
  });

  this.selectedImages = [];
  this.showRateModal = true;
}

closeRateModal() {
  this.showRateModal = false;
  this.selectedTourId = undefined;
  this.selectedImages = [];
}

onFileChange(evt: Event) {
  const input = evt.target as HTMLInputElement;
  if (!input.files) return;
  this.selectedImages = Array.from(input.files);
}

submitReview() {
  if (!this.selectedTourId) return;
  this.submittingReview = true;

  const v = this.reviewForm.value;
  const images = [v.imageUrl1, v.imageUrl2, v.imageUrl3].filter(Boolean);

  const payload = {
    rating: this.reviewForm.value.rating,
    comment: this.reviewForm.value.comment,
    touristName: this.reviewForm.value.touristName,
    touristEmail: this.reviewForm.value.touristEmail,
    visitedAt: this.reviewForm.value.visitedAt,
    commentedAt: this.reviewForm.value.commentedAt,
    images,
  };

  this.tourService.createReview(this.selectedTourId, payload).subscribe({
    next: () => {
      this.submittingReview = false;
      this.closeRateModal();
    },
    error: (err) => {
      console.error('Slanje recenzije neuspešno', err);
      this.submittingReview = false;
      alert('Nismo uspeli da sačuvamo recenziju. Pokušaj ponovo.');
    }
  });
}

  openRatings(tourId: string) {
    this.selectedTourForRatings = tourId;
    this.showRatingsModal = true;
    this.loadRatings(tourId);
  }

  closeRatingsModal() {
    this.showRatingsModal = false;
    this.ratings = [];
    this.ratingsError = '';
    this.selectedTourForRatings = null;
  }

  private loadRatings(tourId: string) { // <-- string
    this.ratingsLoading = true;
    this.ratingsError = '';

    this.tourService.getRatings(tourId).subscribe({
      next: (res) => {
        this.ratings = res || [];
        this.ratingsLoading = false;
      },
      error: () => {
        this.ratingsError = 'Failed to load ratings.';
        this.ratingsLoading = false;
      }
    });
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

      tour.status = this.statusToNumber('Published');
      tour.publishedAt = new Date().toISOString();

      this.tourService.update(tour.id, tour).subscribe({
        next: () => {
          console.log('Tour published successfully');
        },
        error: (err) => {
          console.error('Failed to publish tour', err);
          tour.status = TourStatus.Draft;
        },
      });
    });
  }

  archiveTour(tour: any) {
    tour.status = this.statusToNumber('Archived');
    tour.archivedAt = new Date().toISOString();

    this.tourService.update(tour.id, tour).subscribe({
      next: () => {
        console.log('Tour archived successfully');
      },
      error: (err) => {
        console.error('Failed to archive tour', err);
        tour.status = TourStatus.Published;
      },
    });
  }

  startTour(tourId: string) {
    if (!this.currentUserId) {
      alert('User not loaded—please log in again.');
      return;
    }
    this.tourService.startExecution({ userId: this.currentUserId, tourId }).subscribe({
      next: (te) => {
        this.router.navigate(['/position-sim', tourId], {
          state: { fromStartTour: true, te }
        });
      },
      error: (e) => console.error('Failed to start execution', e),
    });
  }

}
