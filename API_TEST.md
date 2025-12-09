# Soraka Backend API 测试文档

## 基础信息
- 后端地址: `http://localhost:8080`
- API 前缀: `/api`

## 1. 游戏流程相关 API

### 1.1 获取游戏阶段
```http
GET /api/gameflow/phase
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "phase": "Lobby"
  }
}
```

**游戏阶段枚举:**
- `None` - 无
- `Lobby` - 大厅
- `Matchmaking` - 匹配中
- `ReadyCheck` - 准备确认
- `ChampSelect` - 英雄选择
- `GameStart` - 游戏开始
- `InProgress` - 游戏中
- `Reconnect` - 重连
- `WaitingForStats` - 等待结算
- `EndOfGame` - 游戏结束

### 1.2 获取游戏会话信息
```http
GET /api/gameflow/session
```

### 1.3 重新连接游戏
```http
POST /api/gameflow/reconnect
```

### 1.4 获取准备确认状态
```http
GET /api/gameflow/ready-check
```

### 1.5 接受对局
```http
POST /api/gameflow/accept-ready-check
```

---

## 2. 英雄选择相关 API

### 2.1 获取英雄选择会话
```http
GET /api/champ-select/session
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "localPlayerCellId": 0,
    "myTeam": [...],
    "theirTeam": [...],
    "actions": [...]
  }
}
```

### 2.2 选择英雄
```http
POST /api/champ-select/select-champion
Content-Type: application/json

{
  "actionId": 1234567890,
  "championId": 267,
  "completed": true
}
```

**参数说明:**
- `actionId`: 动作ID（从会话中获取）
- `championId`: 英雄ID
- `completed`: 是否完成选择（true=锁定，false=预选）

### 2.3 禁用英雄
```http
POST /api/champ-select/ban-champion
Content-Type: application/json

{
  "actionId": 1234567890,
  "championId": 157,
  "completed": true
}
```

### 2.4 接受英雄交换
```http
POST /api/champ-select/accept-trade/{tradeId}
```

### 2.5 接受楼层交换
```http
POST /api/champ-select/accept-swap/{swapId}
```

### 2.6 备战席交换
```http
POST /api/champ-select/bench-swap/{championId}
```

**使用场景:** 大乱斗模式从备战席交换英雄

### 2.7 获取当前选择的英雄
```http
GET /api/champ-select/current-champion
```

### 2.8 获取皮肤轮播
```http
GET /api/champ-select/skin-carousel
```

### 2.9 选择皮肤和召唤师技能
```http
POST /api/champ-select/select-skin
Content-Type: application/json

{
  "skinId": 267001,
  "spell1Id": 4,
  "spell2Id": 14
}
```

**参数说明:**
- `skinId`: 皮肤ID
- `spell1Id`: 召唤师技能1 ID（可选）
- `spell2Id`: 召唤师技能2 ID（可选）

### 2.10 摇骰子
```http
POST /api/champ-select/reroll
```

**使用场景:** 大乱斗模式重新随机英雄

---

## 3. 自动化功能 API

### 3.1 自动接受对局
```http
POST /api/automation/accept-ready-check
```

**功能:** 检测到准备确认状态时自动接受

### 3.2 自动选择英雄
```http
POST /api/automation/select-champion
Content-Type: application/json

{
  "championId": 267
}
```

**功能:** 自动在当前玩家的选择回合中锁定指定英雄

### 3.3 自动禁用英雄
```http
POST /api/automation/ban-champion
Content-Type: application/json

{
  "championId": 157
}
```

**功能:** 自动在当前玩家的禁用回合中禁用指定英雄

### 3.4 自动接受所有英雄交换
```http
POST /api/automation/accept-trades
```

**功能:** 自动接受所有收到的英雄交换请求

### 3.5 自动接受所有楼层交换
```http
POST /api/automation/accept-swaps
```

**功能:** 自动接受所有收到的楼层交换请求

### 3.6 应用符文页 (OPGG一键设置)
```http
POST /api/automation/apply-rune-page
Content-Type: application/json

{
  "name": "OPGG推荐符文",
  "primaryStyleId": 8100,
  "subStyleId": 8200,
  "selectedPerkIds": [8112, 8143, 8138, 8135, 8009, 8014]
}
```

**参数说明:**
- `name`: 符文页名称
- `primaryStyleId`: 主系ID（8000=精密，8100=主宰，8200=巫术，8300=坚决，8400=启迪）
- `subStyleId`: 副系ID
- `selectedPerkIds`: 6个符文ID的数组

**功能:** 自动删除当前符文页（如果可删除），然后创建新符文页

---

## 4. 召唤师信息 API (已有)

### 4.1 获取当前召唤师排位信息
```http
GET /api/current-summoner/rank-info
```

---

## 5. 客户端信息 API (已有)

### 5.1 获取客户端连接状态
```http
GET /api/client/client-info
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "connected": true,
    "gamePhase": "Lobby",
    "token": "xxx",
    "port": 12345,
    "clientPath": "C:\\Riot Games\\League of Legends\\LeagueClient.exe"
  }
}
```

### 5.2 打开英雄联盟客户端
```http
POST /api/client/open-lol-client
```

---

## 常用英雄 ID

| 英雄名 | ID |
|--------|-----|
| 亚索 | 157 |
| 劫 | 238 |
| 娜美 | 267 |
| 阿卡丽 | 84 |
| 盖伦 | 86 |
| 拉克丝 | 99 |
| EZ | 81 |
| 卡牌 | 4 |
| 瑞兹 | 13 |

---

## 常用召唤师技能 ID

| 技能名 | ID |
|--------|-----|
| 闪现 | 4 |
| 传送 | 12 |
| 点燃 | 14 |
| 治疗 | 7 |
| 屏障 | 21 |
| 虚弱 | 3 |
| 净化 | 1 |
| 惩戒 | 11 |

---

## 符文系 ID

| 符文系 | ID |
|--------|-----|
| 精密 | 8000 |
| 主宰 | 8100 |
| 巫术 | 8200 |
| 坚决 | 8300 |
| 启迪 | 8400 |

---

## 测试流程示例

### 场景1: 自动接受对局并选择英雄

1. 启动后端服务
2. 启动英雄联盟客户端
3. 开始匹配
4. 当匹配成功时，调用自动接受对局API:
   ```bash
   curl -X POST http://localhost:8080/api/automation/accept-ready-check
   ```
5. 进入BP阶段后，调用自动选择英雄API:
   ```bash
   curl -X POST http://localhost:8080/api/automation/select-champion \
     -H "Content-Type: application/json" \
     -d '{"championId": 267}'
   ```

### 场景2: OPGG一键设置符文

进入英雄选择阶段后：
```bash
curl -X POST http://localhost:8080/api/automation/apply-rune-page \
  -H "Content-Type: application/json" \
  -d '{
    "name": "主宰-电刑",
    "primaryStyleId": 8100,
    "subStyleId": 8200,
    "selectedPerkIds": [8112, 8143, 8138, 8135, 8009, 8014]
  }'
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 500 | 服务器内部错误 |
| -1 | LCU客户端未连接 |

---

## 注意事项

1. **LCU 连接状态**: 所有LCU相关API都需要先启动英雄联盟客户端
2. **游戏阶段**: 某些API只在特定游戏阶段可用
   - 英雄选择API只在 `ChampSelect` 阶段可用
   - 准备确认API只在 `ReadyCheck` 阶段可用
3. **请求频率**: 避免过于频繁的API调用，建议间隔至少500ms
4. **自动化功能**: 使用自动化API时，确保当前处于正确的游戏阶段

---

## 下一步待实现功能

- [ ] WebSocket 实时事件推送
- [ ] 对局历史详情查询
- [ ] OPGG 数据集成
- [ ] 大乱斗 Buff 信息
- [ ] 游戏版本数据管理
