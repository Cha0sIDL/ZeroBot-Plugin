package HorseRace

type horse struct {
	horseName       string //马的名字
	playerName      string //QQ名字
	playerUid       int64  //QQ号
	selfBuff        []buff
	delayEvent      []event
	round           int
	location        int
	locationAddMove int
	locationAdd     int
}

type globalGame struct { //全局游戏数据结构
	players      []horse
	round        int //全局的回合数
	start        int //用于控制游戏状态
	time         int64
	raceOnlyKeys []string
	events       []event
}

type buff struct { //buff结构
	buffName    string
	roundStart  int
	roundEnd    int
	moveMin     int
	moveMax     int
	eventInBuff []probabilityEvent
	buffTag     string
}

type event struct { //配置事件的数据结构
	RaceOnlyExist         int                `json:"race_only_exist,omitempty"`
	EventName             string             `json:"event_name,omitempty"`
	Describe              string             `json:"describe,omitempty"`
	Target                int                `json:"target,omitempty"`
	TargetIsBuff          string             `json:"target_is_buff,omitempty"`
	TargetNoBuff          string             `json:"target_no_buff,omitempty"`
	Live                  int                `json:"live,omitempty"`
	Move                  int                `json:"move,omitempty"`
	TrackToLocationExist  int                `json:"track_to_location_exist,omitempty"`
	TrackRandomLocation   int                `json:"track_random_location,omitempty"`
	BuffTimeAdd           int                `json:"buff_time_add,omitempty"`
	DelBuff               string             `json:"del_buff,omitempty"`
	TrackExchangeLocation int                `json:"track_exchange_location,omitempty"`
	RandomEventOnce       []probabilityEvent `json:"random_event_once,omitempty"` //修改了结构
	Die                   int                `json:"die,omitempty"`
	DieName               string             `json:"die_name,omitempty"`
	Away                  int                `json:"away,omitempty"`
	AwayName              string             `json:"away_name,omitempty"`
	Rounds                int                `json:"rounds,omitempty"`
	Name                  string             `json:"name,omitempty"`
	MoveMax               int                `json:"move_max,omitempty"`
	MoveMin               int                `json:"move_min,omitempty"`
	LocateLock            int                `json:"locate_lock,omitempty"`
	Vertigo               int                `json:"vertigo,omitempty"`
	Hiding                int                `json:"hiding,omitempty"`
	OtherBuff             []string           `json:"other_buff,omitempty"`
	RandomEvent           []probabilityEvent `json:"random_event,omitempty"` //修改了结构
	DelayEvent            []otherEvent       `json:"delay_event,omitempty"`
	DelayEventSelf        []otherEvent       `json:"delay_event_self,omitempty"`
	AnotherEvent          *event             `json:"another_event,omitempty"`
	AnotherEventSelf      *event             `json:"another_event_self,omitempty"`
	AddHorse              *addHorse          `json:"add_horse,omitempty"`
	ReplaceHorse          *addHorse          `json:"replace_horse,omitempty"`
}

type addHorse struct {
	Horsename string `json:"horsename,omitempty"`
	Uid       int64  `json:"uid,omitempty"`
	Owner     string `json:"owner,omitempty"`
	Location  int    `json:"location,omitempty"`
}

type otherEvent struct {
	Round int   `json:"round,omitempty"`
	Other event `json:"other"`
}

type probabilityEvent struct {
	Probability int   `json:"probability,omitempty"`
	Other       event `json:"other"`
}
