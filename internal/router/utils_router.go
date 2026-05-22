package router

import (
	"bot_story_generator/internal/models"
	"context"
)

func outboundQueueIndex(userID int64, numWorkers int) int {
	idx := int(userID % int64(numWorkers))
	if idx < 0 {
		idx += numWorkers
	}
	return idx
}

func (r *StoryRouterImpl) createOutboundMessage(ctx context.Context, userID int64, text string, butargs ...models.ButtonArg) {
	queue := r.chans_outbound[outboundQueueIndex(userID, r.numworkers)]
	select {
	case <-r.ctx.Done():
		return
	case queue <- models.NewOutboundMessage(ctx, userID, text, butargs...):
	}
}
func (r *StoryRouterImpl) createEditMessage(ctx context.Context, userID int64, msgID int, text string, butargs ...models.ButtonArg) {
	select {
	case <-r.ctx.Done():
		return
	case r.chan_edit <- models.NewEditMessage(ctx, userID, msgID, text, butargs...):
	}
}
func (r *StoryRouterImpl) createDeleteMessage(ctx context.Context, userID int64, msgID int) {
	select {
	case <-r.ctx.Done():
		return
	case r.chan_delete <- models.NewDeleteMessage(ctx, userID, msgID):
	}
}
func (r *StoryRouterImpl) createInvoiceMessage(ctx context.Context, sub *models.Subscription) {
	select {
	case <-r.ctx.Done():
		return
	case r.chan_bot_invoice <- models.NewInvoiceMessage(ctx, sub):
	}
}
func (r *StoryRouterImpl) createPaymentMessage(pm *models.PaymentData) {
	select {
	case <-r.ctx.Done():
		return
	case r.chan_outbound_payments <- pm:
	}
}
func (r *StoryRouterImpl) cleanUserState(userID int64) {
	r.mux.Lock()
	delete(r.userState, userID)
	r.mux.Unlock()
}

func (r *StoryRouterImpl) checkAdmin(userID int64) bool {
	r.mux.RLock()
	_, ok := r.admins[userID]
	r.mux.RUnlock()
	return ok
}
