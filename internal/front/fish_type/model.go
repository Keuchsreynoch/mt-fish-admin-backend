package fishtype

const (
	ActiveStatusID int16 = 1
)

type FishTypesResponse struct {
	FishTypes []FishInfoResponse `json:"fish_types"`
}

type ManifestRenderType struct {
	ID       int    `json:"id"`
	TypeCode string `json:"type_code"`
}

type ManifestRenderState struct {
	ID             int      `json:"id"`
	RenderFamilyID int      `json:"render_family_id"`
	StateCode      string   `json:"state_code"`
	Prefix         string   `json:"prefix"`
	FrameArray     []string `json:"frame_array"`
	IsDefault      bool     `json:"is_default"`
}

type ManifestRenderFamily struct {
	ID               int                   `json:"id"`
	RenderFamilyName string                `json:"render_family_name"`
	RenderTypeID     int                   `json:"render_type_id"`
	SpineJSONPath    *string               `json:"spine_json_path"`
	SpineAtlasPath   *string               `json:"spine_atlas_path"`
	SpinePNGPath     *string               `json:"spine_png_path"`
	RenderType       *ManifestRenderType   `json:"render_type"`
	RenderStates     []ManifestRenderState `json:"render_states"`
}

type FishInfoResponse struct {
	FishTypeName      string   `json:"fish_type_name"`
	IsBoss            bool     `json:"is_boss"`
	BossName          *string  `json:"boss_name"`
	MinKillOdd        *float64 `json:"min_kill_odd"`
	MaxKillOdd        *float64 `json:"max_kill_odd"`
	BaseSpeed         float64  `json:"base_speed"`
	MissRewardEnabled bool     `json:"miss_reward_enabled"`
	MinMissRewardOdd  *float64 `json:"min_miss_reward_odd"`
	MaxMissRewardOdd  *float64 `json:"max_miss_reward_odd"`
}

type manifestFishRow struct {
	FishTypeName      string   `db:"fish_type_name"`
	IsBoss            bool     `db:"is_boss"`
	BossName          *string  `db:"boss_name"`
	MinKillOdd        *float64 `db:"min_kill_odd"`
	MaxKillOdd        *float64 `db:"max_kill_odd"`
	BaseSpeed         float64  `db:"base_speed"`
	MissRewardEnabled bool     `db:"miss_reward_enabled"`
	MinMissRewardOdd  *float64 `db:"min_miss_reward_odd"`
	MaxMissRewardOdd  *float64 `db:"max_miss_reward_odd"`
}
