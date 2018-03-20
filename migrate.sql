CREATE TABLE IF NOT EXISTS `active_task_instance` (
      `instance_uuid` VARCHAR(255) NOT NULL,
      `task_id` INT NOT NULL,
      `instance_state` INT NOT NULL,
      `reminder_state` INT NULL,
      `reminder_expression` VARCHAR(255) NULL,
      `user_id` INT NOT NULL,
      `notify_to_id` BIGINT NOT NULL,
      `start_at` TIMESTAMP NULL,
      `end_at` TIMESTAMP NULL,
      `wander_protect` INT NOT NULL);

CREATE TABLE IF NOT EXISTS `users` (
    `uuid` VARCHAR(255) NOT NULL,
    `id` INT NOT NULL,
    `user_name` VARCHAR(1024) DEFAULT "" NOT NULL,
    `disp_name` VARCHAR(1024) DEFAULT "" NOT NULL,
    `create_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `update_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
    );

