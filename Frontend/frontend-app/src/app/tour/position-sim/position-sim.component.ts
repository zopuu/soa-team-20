import { AfterViewInit, Component, ElementRef, OnDestroy, ViewChild } from '@angular/core';
import * as L from 'leaflet';
import { TouristLocationService, Coordinates } from '../tourist-location.service';
import { AuthService } from '../../auth/auth.service';
import { ActivatedRoute } from '@angular/router';
import { KeypointService } from '../keypoint.service';
import { KeyPoint } from '../keypoint.model';
import { TourService } from '../tour.service';
import { EMPTY, interval, Subscription, switchMap } from 'rxjs';
type VisitedRow = { id: string; title: string; visitedAt: Date };

@Component({
  selector: 'app-position-sim',
  templateUrl: './position-sim.component.html',
  styleUrls: ['./position-sim.component.css']
})
export class PositionSimComponent implements AfterViewInit, OnDestroy {
  @ViewChild('map', { static: true }) mapEl!: ElementRef<HTMLDivElement>;

  private map?: L.Map;
  private marker?: L.Marker;
  private icon = L.icon({
    iconUrl: 'assets/icons/positionmark.svg',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  private kpIcon = L.icon({ iconUrl: 'assets/icons/blue-pin.png', iconSize: [24, 24], iconAnchor: [12, 24] });
  private kpIconVisited = L.icon({ iconUrl: 'assets/icons/gray-pin.png', iconSize: [24, 24], iconAnchor: [12, 24] });

  private checkpointLayer?: L.LayerGroup;  // holds read-only checkpoint markers
  private routeControl?: any;

  tourId?: string;
  fromStartTour = false;
  //private routeControl?: any;   // routing line connecting checkpoints


  userId?: number;
  selected?: Coordinates;
  statusMsg = '';      // shows inline “Saving… / Saved / Failed”
  hasKnownPosition = false; // for the initial hint
  private proximitySub?: Subscription;
  private visitedIds = new Set<string>();
  visited: VisitedRow[] = [];

  totalKeypoints = 0;
  progressPct = 0;


  toastMsg = '';

  nextTitle: string | undefined;
  distanceMeters: number | undefined;

  constructor(
    private api: TouristLocationService,
    private auth: AuthService,
    private route: ActivatedRoute,     // <-- added
    private keypoints: KeypointService, // <-- added
    private tours: TourService,
  ) {}

  ngAfterViewInit(): void {
    this.auth.whoAmI().subscribe({
      next: (me) => {
        this.userId = me.id;
        this.initMap();
        this.loadExisting();
        // Seed visited from active execution (if any)
        this.tours.getActive(this.userId).subscribe({
          next: (te) => {
            if (te?.keyPointsVisited?.length) {
              // seed both the Set and the sidebar rows
              for (const v of te.keyPointsVisited) {
                const id = String(v.id);
                this.visitedIds.add(id);
                this.visited.unshift({
                  id,
                  title: v.title ?? 'Keypoint',
                  visitedAt: v.visitedAt ? new Date(v.visitedAt) : new Date()
                });
              }
              // keep newest first
              this.visited.sort((a,b) => b.visitedAt.getTime() - a.visitedAt.getTime());
            }
            this.readTourIdAndLoad();
            this.startProximityTicker();
          },
          error: () => { this.readTourIdAndLoad(); this.startProximityTicker(); }
        });
      },
      error: () => { this.statusMsg = 'Cannot determine user identity.'; }
    });
  }

  private startProximityTicker() {
    this.proximitySub = interval(10_000).pipe(
      switchMap(() => {
        if (!this.userId || !this.selected || !this.tourId) return EMPTY;
        return this.tours.checkProximity(this.userId, this.selected);
      })
    ).subscribe({
      next: (res) => {
        // always update “distance to next”
        this.nextTitle = res?.nextKeyPoint?.title ?? undefined;
        this.distanceMeters = typeof res?.distanceMeters === 'number' ? res.distanceMeters : undefined;

        if (res?.reached && res?.justCompletedPoint?.id) {
          this.markVisitedAndRefresh(res.justCompletedPoint);
        }

        if (res?.completedSession) {
          this.toast('Tour completed 🎉');
          this.stopTicker();
        }
      },
      error: (err) => console.error('Proximity ticker failed', err)
    });
  }

  private stopTicker() {
    this.proximitySub?.unsubscribe();
    this.proximitySub = undefined;
  }

  private initMap() {
    this.map = L.map(this.mapEl.nativeElement, { center: [44.7866, 20.4489], zoom: 6 });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Click → set marker → auto-save
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const coords: Coordinates = {
        latitude: Number(e.latlng.lat.toFixed(6)),
        longitude: Number(e.latlng.lng.toFixed(6))
      };
      this.placeMarker(coords);
      this.autoSave(coords);
    });
  }

  private readTourIdAndLoad(): void {
    // param or query param
    const paramId = this.route.snapshot.paramMap.get('tourId');
    const queryId = this.route.snapshot.queryParamMap.get('tourId');

    // navigation state (e.g. this.router.navigate(['/position-sim'], { state: { tourId } }))
    const stateId = (history.state && (history.state.tourId ?? null)) as string | number | null;

    const raw = (paramId ?? queryId ?? (stateId != null ? String(stateId) : null));
    if (!raw) return;

    this.tourId = String(raw);
    this.fromStartTour = true;
    this.loadCheckpoints(this.tourId);
  }

  // private loadCheckpoints(tourId: string): void {
  //   this.keypoints.getByTourSorted(tourId).subscribe({
  //     next: (points) => this.renderCheckpoints(points),
  //     error: (err) => console.error('Failed to load checkpoints for tour', tourId, err)
  //   });
  // }

private renderCheckpoints(points: KeyPoint[]): void {
  if (!this.map) return;

  if (this.checkpointLayer) this.map.removeLayer(this.checkpointLayer);
  this.checkpointLayer = L.layerGroup();

  this.totalKeypoints = points?.length ?? 0;            // ← total for progress
  const latlngs: L.LatLng[] = [];

  for (const kp of points ?? []) {
    const id = String(kp.id);
    if (this.visitedIds.has(id)) continue;              // ← HIDE visited ones

    const { latitude: lat, longitude: lng } = kp.coordinates ?? {};
    if (typeof lat !== 'number' || typeof lng !== 'number') continue;

    latlngs.push(L.latLng(lat, lng));
    L.marker([lat, lng], { draggable: false, icon: this.kpIcon })
      .bindPopup(kp.title ? `Checkpoint: ${kp.title}` : 'Checkpoint')
      .addTo(this.checkpointLayer);
  }

  this.checkpointLayer.addTo(this.map);

  // route line among remaining points
  if (latlngs.length > 1) {
    if (this.routeControl) { this.map!.removeControl(this.routeControl); this.routeControl = undefined; }
    this.routeControl = L.Routing.control({
      waypoints: latlngs,
      router: L.routing.mapbox('pk.eyJ1IjoidmVsam9vMDIiLCJhIjoiY20yaGV5OHU4MDFvZjJrc2Q4aGFzMTduNyJ9.vSQUDO5R83hcw1hj70C-RA', { profile: 'mapbox/walking' }),
      routeWhileDragging: false, show: false, addWaypoints: false, fitSelectedRoutes: true,
      lineOptions: { styles: [{ color: 'blue', weight: 4 }], extendToWaypoints: false, missingRouteTolerance: 0 },
      createMarker: () => null
    } as any).addTo(this.map!);
  }

  const layers: L.Layer[] = this.checkpointLayer.getLayers();
  if (this.marker) layers.push(this.marker);
  if (layers.length > 0) {
    const fg = L.featureGroup(layers as any);
    this.map!.fitBounds(fg.getBounds().pad(0.25));
  }

  // update progress %
  this.progressPct = this.totalKeypoints ? Math.round((this.visitedIds.size / this.totalKeypoints) * 100) : 0;
}


private toast(text: string) {
  this.toastMsg = text;
  setTimeout(() => (this.toastMsg = ''), 2000);
}

private markVisitedAndRefresh(justCompleted: any) {
  const id = String(justCompleted.id);
  if (this.visitedIds.has(id)) return;

  this.visitedIds.add(id);
  this.visited.unshift({ id, title: justCompleted.title ?? 'Keypoint', visitedAt: new Date() });
  this.progressPct = this.totalKeypoints ? Math.round((this.visitedIds.size / this.totalKeypoints) * 100) : 0;
  this.toast(`Visited: ${justCompleted.title ?? 'Keypoint'}`);
  if (this.tourId) this.loadCheckpoints(this.tourId);
}


  private loadExisting() {
    if (!this.userId) return;
    this.api.get(this.userId).subscribe({
      next: (loc) => {
        if (loc && this.map) {
          const { latitude, longitude } = loc.coordinates;
          this.placeMarker({ latitude, longitude });
          this.map.setView([latitude, longitude], 13);
          this.hasKnownPosition = true;
        } else {
          this.hasKnownPosition = false; // show the hint
        }
      },
      error: () => { this.hasKnownPosition = false; }
    });
  }

  private loadCheckpoints(tourId: string): void {
    // ⬇️ Clear previous route line if present
    if (this.routeControl && this.map) {
      this.map.removeControl(this.routeControl);
      this.routeControl = undefined;
    }

    this.keypoints.getByTourSorted(tourId).subscribe({
      next: (points) => this.renderCheckpoints(points),
      error: (err) => console.error('Failed to load checkpoints for tour', tourId, err)
    });
  }

  private placeMarker(coords: Coordinates) {
    if (!this.map) return;
    if (this.marker) this.marker.setLatLng([coords.latitude, coords.longitude]);
    else this.marker = L.marker([coords.latitude, coords.longitude], { icon: this.icon }).addTo(this.map);
    this.selected = coords;
  }

  private autoSave(coords: Coordinates) {
    if (!this.userId) return;
    this.statusMsg = 'Saving…';
    this.api.set(this.userId, coords).subscribe({
      next: () => {
        this.statusMsg = 'Saved.';
        this.hasKnownPosition = true;
        setTimeout(() => (this.statusMsg = ''), 1500);

        this.tours.checkProximity(this.userId!, coords).subscribe({
          next: (res) => {
            this.nextTitle = res?.nextKeyPoint?.title ?? undefined;
            this.distanceMeters = typeof res?.distanceMeters === 'number' ? res.distanceMeters : undefined;

            if (res?.reached && res?.justCompletedPoint?.id) {
              this.markVisitedAndRefresh(res.justCompletedPoint);
            }
            if (res?.completedSession) {
              this.toast('Tour completed 🎉');
              this.stopTicker();
            }
          },
          error: (err) => console.error('Check proximity failed', err)
        });
      },
      error: () => { this.statusMsg = 'Failed to save location.'; }
    });
  }

  ngOnDestroy(): void {
    if (this.proximitySub) this.proximitySub.unsubscribe();
    if (this.map) {
      if (this.routeControl) {
        this.map.removeControl(this.routeControl);
        this.routeControl = undefined;
      }
      this.map.remove();
      this.map = undefined;
    }
  }

  onAbandonTour() {
    if (!this.userId) return;
    const go = confirm('Abandon current tour? This will end the execution.');
    if (!go) return;

    this.tours.abandonExecution(this.userId).subscribe({
      next: () => {
        this.toast('Tour abandoned');
        // clear local state
        this.stopTicker();
        this.visitedIds.clear();
        this.visited = [];
        this.progressPct = 0;
        this.nextTitle = undefined;
        this.distanceMeters = undefined;
        // optionally clear remaining markers
        if (this.checkpointLayer && this.map) {
          this.map.removeLayer(this.checkpointLayer);
          this.checkpointLayer = undefined;
        }
        if (this.routeControl && this.map) {
          this.map.removeControl(this.routeControl);
          this.routeControl = undefined;
        }
      },
      error: (err) => {
        console.error(err);
        this.toast('Failed to abandon tour');
      }
    });
  }

}
