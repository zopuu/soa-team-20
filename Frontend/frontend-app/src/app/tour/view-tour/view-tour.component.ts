import {
  Component,
  AfterViewInit,
  ViewChild,
  ElementRef,
  OnDestroy,
  OnInit,
} from '@angular/core';
import * as L from 'leaflet';
import 'leaflet-routing-machine';
import { ActivatedRoute, Router } from '@angular/router';
import { KeypointService } from '../keypoint.service';
import { TourService } from '../tour.service';
import { MatDialog } from '@angular/material/dialog';
import { KeypointDetailDialogComponent } from '../keypoint-detail-dialog/keypoint-detail-dialog.component';

@Component({
  selector: 'app-view-tour',
  templateUrl: './view-tour.component.html',
  styleUrls: ['./view-tour.component.css'],
})
export class ViewTourComponent implements OnInit, AfterViewInit, OnDestroy {
  private routeControl?: L.Routing.Control;
  @ViewChild('mapContainer', { static: true })
  mapContainer!: ElementRef<HTMLDivElement>;
  private map?: L.Map;
  private markers: L.Marker[] = [];
  private existingIcon: L.Icon = L.icon({
    iconUrl: 'assets/icons/keypointmark.svg',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  tourId?: string;
  tour: any = null;

  constructor(
    private route: ActivatedRoute,
    private keypointService: KeypointService,
    private tourService: TourService,
    private dialog: MatDialog,
    private router: Router
  ) {}

  ngOnInit(): void {
    // prefer snapshot for immediate value, also subscribe to changes
    const sid = this.route.snapshot.paramMap.get('id');
    if (sid) {
      this.tourId = sid;
      this.loadTourDetails();
    }

    this.route.paramMap.subscribe((pm) => {
      const id = pm.get('id');
      if (id && id !== this.tourId) {
        this.tourId = id;
        console.log('Tour ID: ' + this.tourId);
        this.loadTourDetails();
        // if map already initialized, reload keypoints
        if (this.map) this.loadKeypoints();
      }
    });
  }

  private loadTourDetails(): void {
    if (!this.tourId) return;
    console.log('Loading tour details for id=' + this.tourId);
    this.tourService.getById(this.tourId).subscribe({
      next: (t) => (this.tour = t),
      error: (err) => console.error('Failed to load tour', err),
    });
  }

  ngAfterViewInit(): void {
    this.map = L.map(this.mapContainer.nativeElement, {
      center: [20, 0],
      zoom: 2,
    });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);

    setTimeout(() => this.map?.invalidateSize(), 200);
    this.loadKeypoints();
  }

  private loadKeypoints(): void {
    if (!this.tourId || !this.map) return;
    this.keypointService.getByTourSorted(this.tourId).subscribe({
      next: (kps) => {
        this.markers.forEach((m) => m.remove());
        this.markers = [];
        const latlngs: L.LatLng[] = [];
        (kps || []).forEach((kp: any) => {
          const lat = kp.coordinates?.latitude;
          const lng = kp.coordinates?.longitude;
          if (lat == null || lng == null) return;
          latlngs.push(L.latLng(lat, lng));
          const m = L.marker([lat, lng], { icon: this.existingIcon }).addTo(
            this.map!
          );
          m.on('click', () => {
            const ref = this.dialog.open(KeypointDetailDialogComponent, {
              data: kp,
              width: '320px',
            });
            ref.afterClosed().subscribe((res) => {
              if (!res) return;
              if (res.action === 'deleted' || res.action === 'delete') {
                this.loadKeypoints();
              }
            });
          });
          this.markers.push(m);
        });
        // Draw street-connected route using Leaflet Routing Machine
        if (latlngs.length > 1) {
          if (this.routeControl) {
            this.map!.removeControl(this.routeControl);
            this.routeControl = undefined;
          }
          this.routeControl = L.Routing.control({
            waypoints: latlngs,
            router: L.routing.mapbox(
              'pk.eyJ1IjoidmVsam9vMDIiLCJhIjoiY20yaGV5OHU4MDFvZjJrc2Q4aGFzMTduNyJ9.vSQUDO5R83hcw1hj70C-RA',
              { profile: 'mapbox/walking' }
            ),
            show: false,
            lineOptions: {
              styles: [{ color: 'blue', weight: 4 }],
              extendToWaypoints: false,
              missingRouteTolerance: 0,
            },
            createMarker: () => null,
          } as any).addTo(this.map!);
        }
      },
      error: (err) => console.error('failed to load keypoints', err),
    });
  }

  ngOnDestroy(): void {
    if (this.routeControl && this.map) {
      this.map.removeControl(this.routeControl);
      this.routeControl = undefined;
    }
    if (this.map) {
      this.map.remove();
      this.map = undefined;
    }
  }

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

  goToCreateKeypoint() {
    // navigate to route that allows creating keypoint for a tour
    this.router.navigate(['/tours', this.tourId, 'create-keypoint']);
  }
}
