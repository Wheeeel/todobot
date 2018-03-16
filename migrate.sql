CREATE TABLE `active_task_instance` (
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

