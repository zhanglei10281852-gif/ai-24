# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

纽约区域在 2027 年春季切换夏令时那天，次日 00:30 的运行被算进前一天，结果一条在前一天中午仍有额度的计划报了 `daily limit reached`。先不要修改代码，只还原根因；生产代码、测试和配置均保持原样。请列出相邻两个本地业务日各自对应的 UTC 窗口，并把计数链路解释清楚。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/ai-24
- 仓库地址：https://github.com/zhanglei10281852-gif/ai-24.git
- parent SHA：e742f1f2e619d554cbeac9cc1b3c02d296a413cc

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/ai-24.git bug-repro
cd bug-repro
git checkout --detach e742f1f2e619d554cbeac9cc1b3c02d296a413cc
go test ./internal/service -run ^TestDailyLimitSeparatesAdjacentLocalBusinessDays$ -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/service -run ^TestDailyLimitSeparatesAdjacentLocalBusinessDays$ -count=1
--- FAIL: TestDailyLimitSeparatesAdjacentLocalBusinessDays (0.50s)
    annotation_intake_behavior_test.go:59: plan on preceding local business day: data_zone conflict: daily run limit reached
FAIL
FAIL	github.com/zhanglei10281852-gif/ai/internal/service	0.505s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/service -run ^TestDailyLimitSeparatesAdjacentLocalBusinessDays$ -count=1
--- FAIL: TestDailyLimitSeparatesAdjacentLocalBusinessDays (1.44s)
    annotation_intake_behavior_test.go:59: plan on preceding local business day: data_zone conflict: daily run limit reached
FAIL
FAIL	github.com/zhanglei10281852-gif/ai/internal/service	1.632s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

诊断必须定位 internal/domain/data_zone.go 的 FixedBusinessDayWindow，列出测试场景中相邻两个纽约本地业务日各自正确的 UTC 起止边界，并证明春季切换日应为 23 小时；需说明 start.Add(24*time.Hour) 如何把次日 00:30 纳入前一窗口，以及 PlanInferenceRun 的计数如何因此误报限额。目标仓库保持零改动。
