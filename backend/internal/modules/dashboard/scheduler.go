package dashboard

import (
	"context"
	"log"
	"time"
)

func StartScheduler(ctx context.Context, m *Module) {
	log.Println("Starting scheduler...")
	dailyRecommendation(ctx, m.Handler.service)
}

func dailyRecommendation(ctx context.Context, service Service) {
	go func() {
		log.Println("Daily recommendation worker started...")
		for {

			now := time.Now()

			// Starts daily at 6 AM
			next := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				13, 40, 0, 0,
				now.Location(),
			)

			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			sleep := time.Until(next)

			log.Println("Next run at:", next)
			log.Println("Sleeping for:", sleep)

			select {
			case <-time.After(sleep):

				log.Println("Running daily recommendation")

				service.InvalidateRecommendationCache(ctx)

				if _, err := service.GenerateDailyRecommendation(ctx); err != nil {
					log.Println("error:", err)
				}

			case <-ctx.Done():
				log.Println("scheduler stopped (ctx canceled)")
				return
			}
		}
	}()
}
