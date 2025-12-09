# Seraphine 功能实现清单

本文档记录了从 Seraphine 迁移到 Soraka 后端的功能实现情况。

## 已实现的后端功能

### 1. 游戏流程管理 (Gameflow)

#### API 端点
- `GET /api/gameflow/phase` - 获取当前游戏流程阶段
- `GET /api/gameflow/session` - 获取游戏流程会话信息
- `POST /api/gameflow/reconnect` - 重新连接游戏
- `GET /api/gameflow/ready-check` - 获取准备确认状态
- `POST /api/gameflow/accept-ready-check` - 接受对局

#### 实现文件
- 仓库层: `internal/dal/repo/gameflow/gameflow.go`
- 业务层: `internal/biz/gameflow/gameflow.go`
- 服务层: `internal/service/gameflow/gameflow.go`
- 响应结构: `internal/dal/resp/gameflow.go`

---

### 2. 英雄选择 (Champion Select)

#### API 端点
- `GET /api/champ-select/session` - 获取英雄选择会话
- `POST /api/champ-select/select-champion` - 选择英雄
- `POST /api/champ-select/ban-champion` - 禁用英雄
- `POST /api/champ-select/accept-trade/:tradeId` - 接受英雄交换
- `POST /api/champ-select/accept-swap/:swapId` - 接受楼层交换
- `POST /api/champ-select/bench-swap/:championId` - 备战席交换
- `GET /api/champ-select/current-champion` - 获取当前选择的英雄
- `GET /api/champ-select/skin-carousel` - 获取皮肤轮播
- `POST /api/champ-select/select-skin` - 选择皮肤和召唤师技能
- `POST /api/champ-select/reroll` - 摇骰子

#### 实现文件
- 仓库层: `internal/dal/repo/champ_select/champ_select.go`
- 业务层: `internal/biz/champ_select/champ_select.go`
- 服务层: `internal/service/champ_select/champ_select.go`
- 响应结构: `internal/dal/resp/champ_select.go`

---

### 3. 符文管理 (Runes)

#### API 功能
- 获取当前符文页
- 删除符文页
- 创建符文页
- 获取所有符文页

#### 实现文件
- 仓库层: `internal/dal/repo/runes/runes.go`
- 响应结构: `internal/dal/resp/runes.go`

---

### 4. 个性化功能 (Profile)

#### API 功能
- 设置个人主页背景皮肤
- 设置在线状态消息
- 设置显示段位
- 设置在线可用性
- 移除挑战令牌
- 移除威望勋章
- 设置召唤师头像

#### 实现文件
- 仓库层: `internal/dal/repo/profile/profile.go`

---

### 5. 大厅管理 (Lobby)

#### API 功能
- 创建5v5训练模式房间

#### 实现文件
- 仓库层: `internal/dal/repo/lobby/lobby.go`

---

### 6. 观战功能 (Spectate)

#### API 功能
- 通过召唤师名字观战

#### 实现文件
- 仓库层: `internal/dal/repo/spectate/spectate.go`

---

### 7. 自动化功能 (Automation)

#### API 端点
- `POST /api/automation/accept-ready-check` - 自动接受对局
- `POST /api/automation/select-champion` - 自动选择英雄
- `POST /api/automation/ban-champion` - 自动禁用英雄
- `POST /api/automation/accept-trades` - 自动接受所有英雄交换
- `POST /api/automation/accept-swaps` - 自动接受所有楼层交换
- `POST /api/automation/apply-rune-page` - 应用符文页（OPGG一键设置）

#### 实现文件
- 业务层: `internal/biz/automation/automation.go`
- 服务层: `internal/service/automation/automation.go`

---

## 配置文件更新

### config.yaml 新增的 LCU API 端点配置

```yaml
lcu:
  # 游戏流程相关
  gameflow_path: "/lol-gameflow/v1/gameflow-phase"
  gameflow_session_path: "/lol-gameflow/v1/session"
  reconnect_path: "/lol-gameflow/v1/reconnect"
  
  # 英雄选择相关
  champ_select_session_path: "/lol-champ-select/v1/session"
  champ_select_actions_path: "/lol-champ-select/v1/session/actions"
  champ_select_trades_path: "/lol-champ-select/v1/session/trades"
  champ_select_swaps_path: "/lol-champ-select/v1/session/swaps"
  
  # 符文相关
  perks_current_page_path: "/lol-perks/v1/currentpage"
  perks_pages_path: "/lol-perks/v1/pages"
  
  # 个性化相关
  regalia_path: "/lol-regalia/v2/current-summoner/regalia"
  challenges_preferences_path: "/lol-challenges/v1/update-player-preferences"
  
  # 大厅和观战
  lobby_path: "/lol-lobby/v2/lobby"
  spectate_launch_path: "/lol-spectator/v1/spectate/launch"
```

---

## 待实现功能

### 1. LCU WebSocket 事件监听
- [ ] 游戏阶段变化事件
- [ ] 英雄选择变化事件
- [ ] 准备确认事件
- [ ] 当前召唤师信息变化事件

### 2. 战绩查询系统
- [ ] 召唤师战绩查询接口扩展
- [ ] 对局历史详细信息
- [ ] 排位数据分析

### 3. 外部数据集成
- [ ] 大乱斗 Buff 信息 API
- [ ] OPGG 数据集成
  - [ ] 英雄排行榜
  - [ ] 出装加点推荐
  - [ ] 符文推荐

### 4. 前端实现
- [ ] 游戏流程状态显示
- [ ] 英雄选择界面
- [ ] 自动化功能配置界面
- [ ] 个性化设置界面
- [ ] 战绩查询界面

### 5. 游戏版本管理
- [ ] 自动更新队列信息
- [ ] 自动更新游戏模式
- [ ] 自动更新英雄数据
- [ ] 自动更新物品和符文数据

---

## 架构说明

### 分层架构
1. **服务层 (Service)**: HTTP API 端点定义和请求处理
2. **业务逻辑层 (Biz/UseCase)**: 业务规则和流程控制
3. **仓库层 (Repository)**: LCU API 调用和数据访问
4. **数据传输对象 (DTO/Resp)**: 请求和响应结构定义

### 依赖注入
使用 Google Wire 进行依赖注入管理，确保各层之间的解耦。

---

## Seraphine 对应功能清单

### ✅ 已实现
- [x] 自动接受对局
- [x] 自动选择英雄
- [x] 自动禁用英雄
- [x] 自动接受英雄/楼层交换
- [x] 创建5v5训练房间
- [x] 观战功能
- [x] 退出后自动重连
- [x] 修改个人主页背景
- [x] 修改在线状态
- [x] 修改段位显示
- [x] 一键卸下勋章
- [x] 一键卸下头像框（威望勋章）
- [x] 符文页管理（创建、删除）

### 🔄 部分实现
- [~] 战绩查询功能（基础API已有，需增强）
  - [x] 同大区召唤师查询
  - [x] 排位信息查询
  - [ ] 进入BP自动查队友
  - [ ] 进入游戏自动查对手

### ⏳ 待实现
- [ ] OPGG数据显示
  - [ ] 英雄排行
  - [ ] 英雄出装加点
  - [ ] 一键设置符文
- [ ] 大乱斗英雄Buff信息
- [ ] 锁定游戏内设置
- [ ] 热重启客户端
- [ ] 修复客户端DPI问题

---

## 下一步计划

1. **实现 WebSocket 事件监听系统**
   - 监听游戏状态变化
   - 自动触发相应功能

2. **完善战绩查询系统**
   - 扩展对局历史接口
   - 添加数据分析功能

3. **集成外部数据源**
   - 实现OPGG API调用
   - 实现大乱斗Buff数据获取

4. **前端开发**
   - 使用Vue3和Tauri创建桌面应用界面
   - 实现各功能模块的UI

5. **游戏版本管理**
   - 实现自动更新机制
   - 保持与游戏版本同步
