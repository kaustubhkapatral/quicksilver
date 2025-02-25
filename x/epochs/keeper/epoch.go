package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/epochs/types"
)

// GetEpochInfo returns epoch info by identifier.
func (k *Keeper) GetEpochInfo(ctx sdk.Context, identifier string) types.EpochInfo {
	epoch := types.EpochInfo{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEpoch)
	bz := store.Get([]byte(identifier))
	if len(bz) == 0 {
		return epoch
	}

	k.cdc.MustUnmarshal(bz, &epoch)
	return epoch
}

// SetEpochInfo sets epoch info.
func (k *Keeper) SetEpochInfo(ctx sdk.Context, epoch types.EpochInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEpoch)
	bz := k.cdc.MustMarshal(&epoch)
	store.Set([]byte(epoch.Identifier), bz)
}

// DeleteEpochInfo deletes epoch info.
func (k *Keeper) DeleteEpochInfo(ctx sdk.Context, identifier string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEpoch)
	store.Delete([]byte(identifier))
}

// IterateEpochInfo iterates through epochs.
func (k *Keeper) IterateEpochInfo(ctx sdk.Context, fn func(index int64, epochInfo types.EpochInfo) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEpoch)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		epoch := types.EpochInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &epoch)

		stop := fn(i, epoch)

		if stop {
			break
		}
		i++
	}
}

// AllEpochInfos returns every epochInfo in the store.
func (k *Keeper) AllEpochInfos(ctx sdk.Context) []types.EpochInfo {
	var epochs []types.EpochInfo
	k.IterateEpochInfo(ctx, func(_ int64, epochInfo types.EpochInfo) (stop bool) {
		epochs = append(epochs, epochInfo)
		return false
	})
	return epochs
}
