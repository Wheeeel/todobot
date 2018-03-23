使用指南
====

Telegram(目前仅支持 Telegram)机器人 [@Wheeeel_TODObot](https://t.me/Wheeeel_TODObot) 使用指南

## 用法
* **/todo** `/todo Task content##<enroll_count>,Task content##<enroll_count>...`
* **/ping** `/ping`
* **/list** `/list ["done", "all"]` OR `/list`
* **/done** `/done <TaskID>` OR `/done` OR `/donex<TaskID>`
* **/track** `/track ["on", "off"]`
* **/workon** `/workon <TaskID>`
* **/rank** `/rank [count]` OR `/rank`
* **/help** `/help`
* **/del** `/del <TaskID,TaskID...>`
* **/cancel** `/cancel`

## 命令说明

### `/todo` 
添加 TODO 任务，可以批量添加，通过英文逗号","分割。 并且支持设置该 TODO 的参与人数，设置方式为: `/todo blablabla##<count>`
### `/ping` 
检查机器人的可用性并且会通过 ["一言" API](https://hitokoto.cn/api) 返回一个动漫名句
### `/list` 
默认列出当前群组内所有未完成的任务，可以通过增加参数 `all`,`done` 列出所有任务，和所有完成的任务（目前没有分页所以在任务特别多的 Chat 里会刷不出来全部任务消息 QAQ)
### `/done` 
完成任务，如果当前用户正在该群组内 `workon` 某一个任务，则该任务会被标记为完成。如果用户没有在该群组内 `workon` 某一个任务，则会弹出 Reply Keyboard Button 即消息回复按钮，用户可以选择任务 ID 进行完成。也可以通过 `/done <TaskID>` 的形式直接完成某一个任务， 还可以点击 `/list` 中出现的 `/donex<TaskID>` 直接完成某一个任务。
### `/track`
设置机器人跟踪状态，当你在群组中使用机器人时，默认会追踪你的 username & display name, 追踪的信息将可能被公开展示(如 ranklist, 统计信息等)。使用 `/track off` 即可关闭追踪，并且将 username & display name 都设置为 "HIDDEN BY USER" ，再次输入 `/track` 或者 `/track on` 即可开启追踪（不过 userID 是必须记录的啦，不然波特就无法工作了呢）
### `/workon`
开启摸鱼保护模式，可以通过 `/workon <TaskID>` 或者点击 `/todo` 返回结果下面的按钮(限 Telegram)对任务开启摸鱼保护模式，进入该模式之后，如果用户出现在含有该 bot 的群里（包括发消息，发图片等行为），bot 会对摸鱼行为进行提醒(如果任务为 "睡觉", "休息", "sleep" 则会提醒用户去休息)，每 30s 提醒一次，并且每条消息都会被记录为一次 "摸鱼"，无论 bot 是否提醒用户。在用户完成该任务时，会统计用户的工作时间，以及摸鱼次数，有效工作时间等信息，其中工作时间会发送到群组内，有效工作时间和摸鱼次数则通过私聊的形式发送给用户  **使用本功能前请至少私聊过一次 BOT**
### `/rank`
显示完成任务的排行榜，范围为全部使用该bot的用户，在排行榜中将显示用户的 display name 和 完成任务数量，摸鱼次数，如果用户关闭了 BOT 跟踪状态的话，则 display name 显示为 "HIDDEN BY USER"
### `/help`
显示本帮助信息
### `/del`
用于删除任务，用法为 `/del <TaskID>,<TaskID>...`
### `/cancel`
用于取消显示在用户输入框下方的键盘，目前仅有此用途


#### P.S. 对于其他相关问题，如功能添加 意见反馈请提 Issue 同时欢迎 Pull Request ~~
