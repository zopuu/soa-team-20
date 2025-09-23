// cmd/orchestrator/main.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	shoppingpb "github.com/zopuu/soa-team-20/Backend/services/shopping_service/proto/shoppingpb"
)

type Tour struct {
  Id       string  `json:"id"`
  Title    string  `json:"title"` // or "name" depending on FE; your handler returns Title
  Price    float64 `json:"price"`
}

type startReq struct {
  UserId string `json:"userId"`
  TourId string `json:"tourId"`
}

type startResp struct {
  Id string `json:"id"` // full TourExecution payload if you want to forward it
  // ... add fields you want to bubble up to UI
}

func main() {
  tourBase := env("TOUR_BASE", "http://tour-service:8080")          // your tour service HTTP base
  shopAddr := env("SHOPPING_ADDR", "shopping-service:50052")         // shopping gRPC
  httpAddr := env("HTTP_ADDR", ":8085")                              // orchestrator HTTP
  httpClient := &http.Client{Timeout: 10 * time.Second}

  // gRPC client
  cc, err := grpc.Dial(shopAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil { log.Fatalf("dial shopping: %v", err) }
  defer cc.Close()
  shopping := shoppingpb.NewShoppingServiceClient(cc)

  mux := http.NewServeMux()
  mux.HandleFunc("/orchestrations/purchase-start", func(w http.ResponseWriter, r *http.Request) {
    var req startReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
      http.Error(w, "bad request", http.StatusBadRequest); return
    }
    ctx := r.Context()

    // 1) Load tour info (need title+price for cart item)
    tour, err := getTour(ctx, httpClient, tourBase, req.TourId)
    if err != nil {
      http.Error(w, "tour not found: "+err.Error(), http.StatusBadRequest); return
    }

    // 2) Start execution (idempotent on same user+tour as per your Start impl)
    teBody, _ := json.Marshal(map[string]string{"userId": req.UserId, "tourId": req.TourId})
    teResp, err := httpClient.Post(tourBase+"/tours/tour-executions/start", "application/json", bytes.NewReader(teBody))
    if err != nil || teResp.StatusCode >= 400 {
      http.Error(w, "start failed", http.StatusBadRequest); return
    }
    defer teResp.Body.Close()
    var te startResp
    _ = json.NewDecoder(teResp.Body).Decode(&te)

    // 3) Add to cart & checkout (SAGA forward steps)
    if _, err = shopping.AddToCart(ctx, &shoppingpb.AddToCartRequest{
      UserId: req.UserId,
      Item:   &shoppingpb.OrderItem{TourId: req.TourId, Name: firstNonEmpty(tour.Title, "Tour"), Price: tour.Price},
    }); err != nil {
      compensateAbandon(ctx, httpClient, tourBase, req.UserId) // rollback
      http.Error(w, "add to cart failed: "+err.Error(), http.StatusBadRequest); return
    }

    chk, err := shopping.Checkout(ctx, &shoppingpb.CheckoutRequest{UserId: req.UserId})
    if err != nil {
      compensateAbandon(ctx, httpClient, tourBase, req.UserId) // rollback
      http.Error(w, "checkout failed: "+err.Error(), http.StatusBadRequest); return
    }

    // success → reply with both execution & tokens (optional)
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(map[string]any{
      "execution": te,
      "tokens":    chk.Tokens,
    })
  })

  log.Printf("orchestrator listening on %s", httpAddr)
  log.Fatal(http.ListenAndServe(httpAddr, mux))
}

func getTour(ctx context.Context, httpClient *http.Client, base, id string) (*Tour, error) {
  req, _ := http.NewRequestWithContext(ctx, http.MethodGet, base+"/tours/"+id, nil)
  resp, err := httpClient.Do(req); if err != nil { return nil, err }
  defer resp.Body.Close()
  if resp.StatusCode >= 400 { return nil, errors.New(resp.Status) }
  var t Tour; if err := json.NewDecoder(resp.Body).Decode(&t); err != nil { return nil, err }
  return &t, nil
}

func compensateAbandon(ctx context.Context, httpClient *http.Client, base, userId string) {
  body, _ := json.Marshal(map[string]string{"userId": userId})
  _, _ = httpClient.Post(base+"/tours/tour-executions/abandon", "application/json", bytes.NewReader(body))
  // Abandon endpoint exists and flips status; safer than full delete. :contentReference[oaicite:15]{index=15}
}

func env(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
func firstNonEmpty(a, b string) string { if a != "" { return a }; return b }
