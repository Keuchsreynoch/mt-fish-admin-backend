local function init_user_spaces()
  local userSpace = box.space.users_space
  local seq_user_id = box.sequence.seq_user_id

  if userSpace == nil then
    if seq_user_id ~= nil then
      seq_user_id:drop()
    end

    box.schema.sequence.create('seq_user_id', { start = 1 })

    local format = {
      { "id",             "unsigned" },
      { "user_uuid",      "uuid",     is_nullable = false },

      { "first_name",     "string",   is_nullable = false },
      { "last_name",      "string",   is_nullable = false },
      { "user_name",      "string",   is_nullable = false },
      { "password",       "string",   is_nullable = false },
      { "email",          "string",   is_nullable = false },
      { "role_id",        "number",   is_nullable = false },
      { "status",         "boolean",  is_nullable = false },
      { "login_session",  "string",   is_nullable = true },
      { "profile_photo",  "string",   is_nullable = true },
      { "user_alias",     "string",   is_nullable = true },
      { "phone_number",   "string",   is_nullable = true },
      { "user_avatar_id", "number",   is_nullable = true },
      { "commission",     "decimal",  is_nullable = true,  default = decimal.new(0.00) },

      { "status_id",      "number",   is_nullable = false, default = 1 },
      { "order",          "number",   is_nullable = true,  default = 1 },
      { "created_by",     "number",   is_nullable = false },
      { "created_at",     "datetime", is_nullable = false },
      { "updated_by",     "number",   is_nullable = true },
      { "updated_at",     "datetime", is_nullable = true },
      { "deleted_by",     "number",   is_nullable = true },
      { "deleted_at",     "datetime", is_nullable = true },
    }

    userSpace = box.schema.create_space('users_space', { format = format, id = 1001 })

    -- Create index
    userSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_user_id',
      if_not_exists = true,
    })
    userSpace:create_index('user_uuid',
      {
        parts = { { 'user_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      })

    userSpace:create_index('role_id',
    {
      parts = { { 'role_id', 'number' } },
      if_not_exists = true,
      unique = false
    })


    box.sequence.seq_user_id:set(0)
  end

  if userSpace:len() == 0 then
    -- Insert record
    userSpace:auto_increment {
      uuid.fromstr('c5b66b62-2cb0-4a2e-b704-1da97d8ed10d'),
      'Supper',
      'Admin',
      'ADMIN',
      '123',
      'admin@gmail.com',
      1,
      true,
      'bdeb581454a4441784be1e355faeab63',
      'user1.png',
      'KM001',
      '010123123',
      nil,
      decimal.new(0.00),
      1,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }
    userSpace:auto_increment {
      uuid.fromstr('83751b48-68f3-4805-a7bd-60ab8311936d'),
      'IT',
      'Developer',
      'IT',
      '12e!!121#',
      'it@gmail.com',
      1,
      true,
      'bdeb581454a4441784be1e355faeab57',
      'user2.png',
      'KM002',
      '430123123',
      nil,
      decimal.new(0.00),
      1,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }
  end
end

local function init_users_roles_space()
  local UsersRolesSpace = box.space.users_roles_space
  local seq_user_role_id = box.sequence.seq_user_role_id

  if UsersRolesSpace == nil then
    if seq_user_role_id ~= nil then
      seq_user_role_id:drop()
    end

    box.schema.sequence.create('seq_user_role_id', { start = 1 })

    local format = {
      { "id",             "unsigned" },
      { "user_role_uuid", "uuid",     is_nullable = false },

      { "user_role_name", "string",   is_nullable = false },
      { "user_role_desc", "string",   is_nullable = false },
      { "status",         "boolean",  is_nullable = false },
      { "order",          "number",   is_nullable = true, default = 1 },
      { "created_by",     "number",   is_nullable = false },
      { "created_at",     "datetime", is_nullable = false },
      { "updated_by",     "number",   is_nullable = true },
      { "updated_at",     "datetime", is_nullable = true },
      { "deleted_by",     "number",   is_nullable = true },
      { "deleted_at",     "datetime", is_nullable = true },
    }

    UsersRolesSpace = box.schema.create_space('users_roles_space', { format = format, id = 1002 })

    -- Create index
    UsersRolesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_user_role_id',
      if_not_exists = true,
    })
    UsersRolesSpace:create_index('user_role_uuid',
      {
        parts = { { 'user_role_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      })
    box.sequence.seq_user_role_id:set(0)
  end

  if UsersRolesSpace:len() == 0 then
    -- Insert record
    UsersRolesSpace:auto_increment {
      uuid.fromstr('9a6f17b3-f2d1-4df4-8ade-d1c8fbebdb97'),
      'admin',
      'Role Admin',
      true,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }
    -- Insert record
    UsersRolesSpace:auto_increment {
      uuid.fromstr('01918ff3-57fd-7c09-93a3-7077087550b0'),
      'moderator',
      'Role Moderator',
      true,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }
    UsersRolesSpace:auto_increment {
      uuid.fromstr('01918fdb-0d16-7b77-b0c6-dc681aada863'),
      'operator',
      'Role Operator',
      true,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }
  end
end

local function init_users_audits_spaces()
  local UsersAuditsSpace = box.space.users_audits_space
  local seq_user_audit_id = box.sequence.seq_user_audit_id

  if UsersAuditsSpace == nil then
    if seq_user_audit_id ~= nil then
      seq_user_audit_id:drop()
    end

    box.schema.sequence.create('seq_user_audit_id', { start = 1 })

    local format = {
      { "id",                 "unsigned" },
      { "user_audit_uuid",    "uuid",     is_nullable = false },

      { "user_id",            "number",   is_nullable = true, default = 0 },
      { "user_audit_context", "string",   is_nullable = false },
      { "user_audit_desc",    "string",   is_nullable = false },
      { "audit_type_id",      "number",   is_nullable = false },
      { "user_agent",         "string",   is_nullable = false },
      { "operator",           "string",   is_nullable = false },
      { "ip",                 "string",   is_nullable = false },
      { "status_id",          "number",   is_nullable = false },
      { "order",              "number",   is_nullable = true, default = 1 },
      { "created_by",         "number",   is_nullable = false },
      { "created_at",         "datetime", is_nullable = false },
      { "updated_by",         "number",   is_nullable = true },
      { "updated_at",         "datetime", is_nullable = true },
      { "deleted_by",         "number",   is_nullable = true },
      { "deleted_at",         "datetime", is_nullable = true },
    }

    UsersAuditsSpace = box.schema.create_space('users_audits_space', { format = format, id = 1003 })

    -- Create index
    UsersAuditsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_user_audit_id',
      if_not_exists = true,
    })
    UsersAuditsSpace:create_index('user_audit_uuid',
      {
        parts = { { 'user_audit_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      })
    box.sequence.seq_user_audit_id:set(0)
  end
end

local function init_admin_admin_menus_space()
  local adminMenusSpace = box.space.admin_menus_space
  local seq_admin_menu_id = box.sequence.seq_admin_menu_id

  if adminMenusSpace == nil then
    if seq_admin_menu_id ~= nil then
      seq_admin_menu_id:drop()
    end

    box.schema.sequence.create('seq_admin_menu_id', { start = 1 })

    local format = {
      { "id",         "unsigned" },
      { "menu_uuid",  "uuid",     is_nullable = true },
      { "name",       "string",   is_nullable = false },
      { "icon",       "string",   is_nullable = false },
      { "path",       "string",   is_nullable = true },
      { "parent_id",  "number",   is_nullable = false },
      { "status_id",  "number",   is_nullable = false },
      { "order",      "number",   is_nullable = true, default = 1 },
      { "created_by", "number",   is_nullable = false },
      { "created_at", "datetime", is_nullable = false },
      { "updated_by", "number",   is_nullable = true },
      { "updated_at", "datetime", is_nullable = true },
      { "deleted_by", "number",   is_nullable = true },
      { "deleted_at", "datetime", is_nullable = true },
    }

    adminMenusSpace = box.schema.create_space('admin_menus_space', { format = format, id = 1004 })

    adminMenusSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_admin_menu_id',
      if_not_exists = true
    })
  end

  if adminMenusSpace:len() == 0 then
    adminMenusSpace:auto_increment { uuid.fromstr('d12242dd-be25-491c-9355-95a81bfb0581'), 'Dashboard', '/src/assets/images/sidebar/dashboard.svg', '', 0, 1, 1, 1, datetime.parse('2024-08-22T15:54:10.301612345Z'), nil, nil, nil, nil }
    adminMenusSpace:auto_increment { uuid.fromstr('a22c01a2-c6c9-49a0-8de9-1e61e4f21c99'), 'Results', '/src/assets/images/sidebar/results.svg', '', 0, 1, 1, 1, datetime.parse('2024-08-22T15:54:20.301612345Z'), nil, nil, nil, nil }
    adminMenusSpace:auto_increment { uuid.fromstr('06517a8f-eb10-4d78-8f43-7218e814e1fb'), 'Users', '/src/assets/images/sidebar/users.svg', '', 0, 1, 1, 1, datetime.parse('2024-08-22T15:55:30.301612345Z'), nil, nil, nil, nil }
    adminMenusSpace:auto_increment { uuid.fromstr('8fce1c99-1853-4bad-bc43-73c275bedba1'), 'Players', '/src/assets/images/sidebar/players.svg', '', 0, 1, 1, 1, datetime.parse('2024-08-22T15:56:40.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_users_menus_spaces()
  local UsersMenusSpace = box.space.users_menus_spaces
  local seq_user_menu_id = box.sequence.seq_user_menu_id

  if UsersMenusSpace == nil then
    if seq_user_menu_id ~= nil then
      seq_user_menu_id:drop()
    end

    box.schema.sequence.create('seq_user_menu_id', { start = 1 })

    local format = {
      { "id",             "unsigned" },
      { "user_menu_uuid", "uuid",     is_nullable = false },

      { "user_id",        "number",   is_nullable = false },
      { "admin_menu_id",  "number",   is_nullable = false },
      { "status_id",      "number",   is_nullable = false },
      { "order",          "number",   is_nullable = true, default = 1 },
      { "created_by",     "number",   is_nullable = false },
      { "created_at",     "datetime", is_nullable = false },
      { "updated_by",     "number",   is_nullable = true },
      { "updated_at",     "datetime", is_nullable = true },
      { "deleted_by",     "number",   is_nullable = true },
      { "deleted_at",     "datetime", is_nullable = true },
    }

    UsersMenusSpace = box.schema.create_space('users_menus_spaces', { format = format, id = 1005 })

    -- Create index
    UsersMenusSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_user_menu_id',
      if_not_exists = true,
    })
    box.sequence.seq_user_menu_id:set(0)
  end
end

local function init_currencies_space()
  local currenciesSpace = box.space.currencies_space
  local seq_currency_id = box.sequence.seq_currency_id

  if currenciesSpace == nil then
    if seq_currency_id ~= nil then
      seq_currency_id:drop()
    end

    box.schema.sequence.create('seq_currency_id', { start = 1 })

    local format = {
      { "id",            "unsigned" },
      { "currency_uuid", "uuid",     is_nullable = true },
      { "currency_code", "string",   is_nullable = false },
      { "status_id",     "number",   is_nullable = true,  default = 1 },
      { "order",         "number",   is_nullable = false, default = 1 },

      { "created_by",    "number",   is_nullable = false },
      { "created_at",    "datetime", is_nullable = false },
      { "updated_by",    "number",   is_nullable = true },
      { "updated_at",    "datetime", is_nullable = true },
      { "deleted_by",    "number",   is_nullable = true },
      { "deleted_at",    "datetime", is_nullable = true },
    }

    currenciesSpace = box.schema.create_space('currencies_space', { format = format, id = 1006 })

    currenciesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_currency_id',
      if_not_exists = true
    })
    currenciesSpace:create_index('currency_uuid',
      {
        parts = { { 'currency_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if currenciesSpace:len() == 0 then
    -- Insert default records
    currenciesSpace:auto_increment { uuid.fromstr('45828e21-f376-44c8-9dc4-54fc43057c63'), 'KHR', 1, 1, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('d13fa303-9e40-477f-87b1-9192a1707cde'), 'USD', 1, 2, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('cb887d96-ea8f-4d57-a9f5-62dd7aa0ff89'), 'VND', 1, 3, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('d60f2459-b182-47d5-b5db-2a461c4a9227'), 'THB', 1, 4, 1, datetime.parse('2024-08-14T15:48:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('b2dc0c78-e0e0-4fe7-b23d-0c39fc14ed59'), 'CNY', 1, 5, 1, datetime.parse('2024-08-14T15:48:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('6c258b12-d381-4194-8408-c35caa98030b'), 'EUR', 1, 6, 1, datetime.parse('2024-08-14T15:48:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('ea856616-7e10-4095-8c08-7fc1158f7a5f'), 'AUD', 1, 7, 1, datetime.parse('2024-08-14T15:48:30.301612345Z'), nil, nil, nil, nil }
    currenciesSpace:auto_increment { uuid.fromstr('5b158450-e0ee-43e8-be0d-b666f7c45340'), 'CAD', 1, 8, 1, datetime.parse('2024-08-14T15:48:30.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_languages_space()
  local languagesSpace = box.space.languages_space
  local seq_language_id = box.sequence.seq_language_id

  if languagesSpace == nil then
    if seq_language_id ~= nil then
      seq_language_id:drop()
    end

    box.schema.sequence.create('seq_language_id', { start = 1 })

    local format = {
      { "id",            "unsigned" },
      { "language_uuid", "uuid",     is_nullable = true },

      { "country_id",    "number",   is_nullable = false },
      { "language_name", "string",   is_nullable = false },
      { "language_code", "string",   is_nullable = false },
      { "img",           "string",   is_nullable = true },

      { "status_id",     "number",   is_nullable = true,  default = 1 },
      { "order",         "number",   is_nullable = false, default = 1 },
      { "created_by",    "number",   is_nullable = false },
      { "created_at",    "datetime", is_nullable = false },
      { "updated_by",    "number",   is_nullable = true },
      { "updated_at",    "datetime", is_nullable = true },
      { "deleted_by",    "number",   is_nullable = true },
      { "deleted_at",    "datetime", is_nullable = true },
    }

    languagesSpace = box.schema.create_space('languages_space', { format = format, id = 1007 })

    languagesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_language_id',
      if_not_exists = true
    })
    languagesSpace:create_index('language_uuid',
      {
        parts = { { 'language_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if languagesSpace:len() == 0 then
    languagesSpace:auto_increment { uuid.fromstr('e8fc66fd-f42f-4e1a-b6e5-254726acaf40'), 1, 'ភាសាខ្មែរ', 'kh', 'kh.png', 1, 1, 1, datetime.parse('2024-08-14T15:23:10.301612345Z'), nil, nil, nil, nil }
    languagesSpace:auto_increment { uuid.fromstr('77c1c92b-fd4b-482f-9dba-882f82f35bec'), 2, 'English', 'en', 'en.png', 1, 1, 1, datetime.parse('2024-08-14T15:24:20.301612345Z'), nil, nil, nil, nil }
    languagesSpace:auto_increment { uuid.fromstr('f3231590-fd7d-4176-8c0d-55b95efc2c2c'), 3, 'ภาษาไทย', 'th', 'th.png', 1, 1, 1, datetime.parse('2024-08-14T15:25:30.301612345Z'), nil, nil, nil, nil }
    languagesSpace:auto_increment { uuid.fromstr('3d6012ec-0ceb-4094-aedc-5dc5d992908f'), 4, 'tiếng Việt', 'vn', 'vn.png', 1, 1, 1, datetime.parse('2024-08-14T15:26:40.301612345Z'), nil, nil, nil, nil }
    languagesSpace:auto_increment { uuid.fromstr('cfcc8a1a-24be-4423-a428-49d8cdc941bf'), 5, '中文', 'zh', 'cn.png', 1, 1, 1, datetime.parse('2024-08-29T11:27:40.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_exchange_rates_space()
  -- Purpose: Bi-directional exchange rates method.
  -- Eg: From USD to EUR, rate: 0.85 Eg: 20 = 20 * 0.85. Thus, From EUR to USD = 20/0.85
  local exchange_rates_space = box.space.exchange_rates_space
  local seq_exchange_rate_id = box.sequence.seq_exchange_rate_id
  if exchange_rates_space == nil then
    if seq_exchange_rate_id ~= nil then
      seq_exchange_rate_id:drop()
    end

    box.schema.sequence.create('seq_exchange_rate_id', { start = 1 })

    local format = {
      { "id",                 "unsigned" },
      { "exchange_rate_uuid", "uuid",     is_nullable = true },

      { "from_currency_id",   "number",   is_nullable = false },
      { "to_currency_id",     "number",   is_nullable = false },
      { "rate",               "decimal",  is_nullable = false, default = decimal.new(0.000) },

      { "status_id",          "number",   is_nullable = false, default = 1 }, -- 0: Inactive, 1: Active, 2: Deleted for Simplicity. Otherwise, create new fields(effective_from, expiry_date)
      { "order",              "number",   is_nullable = false, default = 1 },
      { "created_by",         "number",   is_nullable = false },
      { "created_at",         "datetime", is_nullable = false },
      { "updated_by",         "number",   is_nullable = true },
      { "updated_at",         "datetime", is_nullable = true },
      { "deleted_by",         "number",   is_nullable = true },
      { "deleted_at",         "datetime", is_nullable = true },
    }

    exchange_rates_space = box.schema.create_space('exchange_rates_space', { format = format, id = 1008 })

    exchange_rates_space:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_exchange_rate_id',
      if_not_exists = true
    })

    box.sequence.seq_exchange_rate_id:set(0)
  end
end

local function init_currencies_defaults_rates_space()
  local currenciesDefaultRateSpace = box.space.currencies_defaults_rates_space
  local seq_currency_default_rate_id = box.sequence.seq_currency_default_rate_id

  if currenciesDefaultRateSpace == nil then
    if seq_currency_default_rate_id ~= nil then
      seq_currency_default_rate_id:drop()
    end

    box.schema.sequence.create('seq_currency_default_rate_id', { start = 1 })

    local format = {
      { "id",                         "unsigned" },
      { "currency_default_rate_uuid", "uuid",     is_nullable = true },
      { "currency_id",                "number",   is_nullable = false },
      { "rate",                       "decimal",  is_nullable = false },

      { "status_id",                  "number",   is_nullable = true,  default = 1 },
      { "order",                      "number",   is_nullable = false, default = 1 },
      { "created_by",                 "number",   is_nullable = false },
      { "created_at",                 "datetime", is_nullable = false },
      { "updated_by",                 "number",   is_nullable = true },
      { "updated_at",                 "datetime", is_nullable = true },
      { "deleted_by",                 "number",   is_nullable = true },
      { "deleted_at",                 "datetime", is_nullable = true },
    }

    currenciesDefaultRateSpace = box.schema.create_space('currencies_defaults_rates_space',
      { format = format, id = 1009 })

    currenciesDefaultRateSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_currency_default_rate_id',
      if_not_exists = true
    })
    currenciesDefaultRateSpace:create_index('status_id',
      {
        parts = { { 'status_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
    currenciesDefaultRateSpace:create_index('currency_id',
      {
        parts = { { 'currency_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

    currenciesDefaultRateSpace:create_index('rate',
      {
        parts = { { 'rate', 'decimal' } },
        if_not_exists = true,
        unique = false
      }
    )

  end

  if currenciesDefaultRateSpace:len() == 0 then
    -- Insert default records
    currenciesDefaultRateSpace:auto_increment { uuid.fromstr('a2a27ec1-aa9c-41b8-8644-3f330c03cdcb'), 1, decimal.new(1), 1, 1, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesDefaultRateSpace:auto_increment { uuid.fromstr('4ce2f032-8501-40dd-a91d-ee12537ca523'), 2, decimal.new(4000), 1, 1, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesDefaultRateSpace:auto_increment { uuid.fromstr('7f9d9b7e-f9b8-4b71-ae70-849c9737f5e3'), 3, decimal.new(0.150), 1, 1, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesDefaultRateSpace:auto_increment { uuid.fromstr('d7b3dab6-6e92-4671-8313-f7fbeac91215'), 4, decimal.new(100), 1, 1, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    currenciesDefaultRateSpace:auto_increment { uuid.fromstr('0bbfa220-e651-4a9e-9977-8c6d7ec6cc79'), 5, decimal.new(500), 1, 1, 1, datetime.parse('2024-08-02T16:05:30.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_channels_space()
  local channelSpace = box.space.channels_space
  local seq_channel_id = box.sequence.seq_channel_id

  if channelSpace == nil then
    if seq_channel_id ~= nil then
      seq_channel_id:drop()
    end

    box.schema.sequence.create('seq_channel_id', { start = 1 })

    local format = {
      { "id",           "unsigned" },
      { "channel_uuid", "uuid",     is_nullable = true },

      { "channel_name", "string",   is_nullable = false },
      { "channel_logo", "string",   is_nullable = true },
      { "stream_one",   "string",   is_nullable = false },
      { "stream_two",   "string",   is_nullable = false },
      { "status_id",    "number",   is_nullable = false, default = 1 },
      { "order",        "number",   is_nullable = false, default = 1 },

      { "created_by",   "number",   is_nullable = false },
      { "created_at",   "datetime", is_nullable = false },
      { "updated_by",   "number",   is_nullable = true },
      { "updated_at",   "datetime", is_nullable = true },
      { "deleted_by",   "number",   is_nullable = true },
      { "deleted_at",   "datetime", is_nullable = true },
    }

    channelSpace = box.schema.create_space('channels_space', { format = format, id = 1010 })

    channelSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_channel_id',
      if_not_exists = true
    })

    channelSpace:create_index('channel_uuid', {
      parts = { { 'channel_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })
  end

  if channelSpace:len() == 0 then
      -- Insert default records for themed channels
      channelSpace:auto_increment {
        uuid.fromstr('7f8d1b9c-3a2e-4f5d-b6c7-8a9d0e1f2b3c'),
        'Tom&Jerry',
        'tomjerry.png',
        'https://live.realtimestreaming.live:5443/LiveApp/play.html?name=TJ-CH1',
        'https://live.realtimestreaming.live:5443/LiveApp/play.html?name=TJ-CH1-ALT',
        1, 1, 1,
        datetime.parse('2024-08-30T09:55:30.301612345Z'),
        nil, nil, nil, nil
      }

      channelSpace:auto_increment {
        uuid.fromstr('a1b2c3d4-5e6f-4a7b-8c9d-0e1f2a3b4c5d'),
        'Spike&Toodles',
        'spiketoodles.png',
        'https://live.realtimestreaming.live:5443/LiveApp/play.html?name=SP-CH2',
        'https://live.realtimestreaming.live:5443/LiveApp/play.html?name=SP-CH2-ALT',
        1, 2, 1,
        datetime.parse('2024-08-30T09:55:30.301612345Z'),
        nil, nil, nil, nil
      }

      channelSpace:auto_increment {
        uuid.fromstr('9e8d7c6b-5a4f-4e3d-2c1b-0a9f8e7d6c5b'),
        'Butch&Nibbles',
        'butchnibbles.png',
        'https://live.realtimestreaming.live:5443/LiveApp/play.html?name=BT-CH3',
        'https://live.realtimestreaming.live:5443/LiveApp/play.html?name=BT-CH3-ALT',
        1, 3, 1,
        datetime.parse('2024-08-30T09:55:30.301612345Z'),
        nil, nil, nil, nil
      }
    end
end

local function init_rounds_space()
  local space_name = 'rounds_space'
  local seq_name = 'seq_round_id'

  if box.space[space_name] ~= nil then
    return
  end

  if box.sequence[seq_name] ~= nil then
    box.sequence[seq_name]:drop()
  end

  box.schema.sequence.create(seq_name, { start = 1 })

  local format = {
    { "id",              "unsigned" },
    { "round_uuid",      "uuid",     is_nullable = true },

    { "room_id",         "number",   is_nullable = false },
    { "round_no",        "string",   is_nullable = false },

    { "start_time",      "datetime", is_nullable = true },
    { "end_time",        "datetime", is_nullable = true },
    { "next_start_time", "datetime", is_nullable = true },

    { "status_id",       "number",   is_nullable = false, default = 1 },
    -- 1:init, 2:betting, 3:dealing, 4:reveal, 5:result, 6:end

    { "issue_by_id",     "number",   is_nullable = false },

    { "created_by",      "number",   is_nullable = false },
    { "created_at",      "datetime", is_nullable = false },
    { "updated_by",      "number",   is_nullable = true },
    { "updated_at",      "datetime", is_nullable = true },
    { "deleted_by",      "number",   is_nullable = true },
    { "deleted_at",      "datetime", is_nullable = true },
  }

  local s = box.schema.create_space(space_name, {
    id = 1011,
    format = format,
  })

  s:create_index('pk', {
    parts = { { 'id', 'unsigned' } },
    sequence = seq_name,
  })

  s:create_index('round_uuid', {
    parts = { { 'round_uuid', 'uuid' } },
    unique = true,
  })

  s:create_index('room_id', {
    parts = { { 'room_id', 'number' } },
     unique = false,
    
  })

  s:create_index('status_id', {
    parts = { { 'status_id', 'number' } },
     unique = false,
  })

  box.sequence[seq_name]:set(0)
end


local function init_results_space()
  local space_name = 'results_space'
  local seq_name = 'seq_result_id'

  if box.space[space_name] ~= nil then
    return
  end

  if box.sequence[seq_name] ~= nil then
    box.sequence[seq_name]:drop()
  end

  box.schema.sequence.create(seq_name, { start = 1 })

  local format = {
    { "id",             "unsigned" },
    { "result_uuid",    "uuid",     is_nullable = true },

    { "round_id",       "number",   is_nullable = false },
    { "room_id",        "number",   is_nullable = false },

    -- Banker cards & result
    { "banker_cards",   "array",    is_nullable = false }, -- msgpack/json cards
    { "banker_score",   "number",   is_nullable = false },
    { "banker_is_pok",  "boolean",  is_nullable = false },
    { "banker_deng",    "number",   is_nullable = false }, -- multiplier

    { "status_id",      "number",   is_nullable = false, default = 1 },
    { "issue_by_id",    "number",   is_nullable = false },

    { "created_by",     "number",   is_nullable = false },
    { "created_at",     "datetime", is_nullable = false },
    { "updated_by",     "number",   is_nullable = true },
    { "updated_at",     "datetime", is_nullable = true },
    { "deleted_by",     "number",   is_nullable = true },
    { "deleted_at",     "datetime", is_nullable = true },
  }

  local s = box.schema.create_space(space_name, {
    id = 1012,
    format = format,
  })

  s:create_index('id', {
    parts = { { 'id', 'unsigned' } },
    sequence = seq_name,
  })

  s:create_index('result_uuid', {
    parts = { { 'result_uuid', 'uuid' } },
    unique = true,
  })

  s:create_index('round_id', {
    parts = { { 'round_id', 'number' } },
  })

  s:create_index('room_id', {
    parts = { { 'room_id', 'number' } },
     unique = false,
  })

  box.sequence[seq_name]:set(0)
end


local function init_cocks_space()
  local cockSpace = box.space.cocks_space
  local seq_cock_id = box.sequence.seq_cock_id

  if cockSpace == nil then
    if seq_cock_id ~= nil then
      seq_cock_id:drop()
    end

    box.schema.sequence.create('seq_cock_id', { start = 1 })

    local format = {
      { "id",         "unsigned" },
      { "cock_uuid",  "uuid",     is_nullable = true },

      { "channel_id", "number",   is_nullable = false, default = 1 },
      { "cock_name",  "string",   is_nullable = false },
      { "status_id",  "number",   is_nullable = false, default = 1 },
      { "order",      "number",   is_nullable = false, default = 1 },

      { "created_by", "number",   is_nullable = false },
      { "created_at", "datetime", is_nullable = false },
      { "updated_by", "number",   is_nullable = true },
      { "updated_at", "datetime", is_nullable = true },
      { "deleted_by", "number",   is_nullable = true },
      { "deleted_at", "datetime", is_nullable = true },
    }

    cockSpace = box.schema.create_space('cocks_space', { format = format, id = 1013 })

    cockSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_cock_id',
      if_not_exists = true
    })

    cockSpace:create_index('cock_uuid', {
      parts = { { 'cock_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })
  end

  if cockSpace:len() == 0 then
    -- Channel 1 (Main characters)
    cockSpace:auto_increment { uuid.fromstr('d4cde0b8-8f4b-4a80-b8ae-6d4e7e5d73b4'), 1, 'Tom', 1, 1, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('f9b82c8b-f3e3-4a45-bf24-5b4b77a47d2a'), 1, 'Jerry', 1, 2, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('3b88d3e5-0bc5-4a38-98aa-88bede8e5a9b'), 1, 'Draw', 1, 3, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('c7c7a8e4-ada0-4790-b00f-486ab6b6e7c9'), 1, 'Cancel', 1, 4, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }

    -- Channel 2 (Supporting characters)
    cockSpace:auto_increment { uuid.fromstr('e1c4e0cf-8c30-4f43-b1ec-1e6f2e7e21b5'), 2, 'Spike', 1, 5, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('c9a7c2b4-5a6e-4a3f-9010-02c9f8f4b98e'), 2, 'Toodles', 1, 6, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('b3e61896-16c5-4c4f-874c-dc78d897c1de'), 2, 'Draw', 1, 7, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('a2b8d2e6-3f4e-4e77-90e1-3c0f9d4b8a25'), 2, 'Cancel', 1, 8, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }

    -- Channel 3 (Recurring / other characters)
    cockSpace:auto_increment { uuid.fromstr('ac7cddc0-f4a7-47b6-a7b8-f2d5f3a8b5c9'), 3, 'Butch', 1, 9, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('eab3f2b2-b9d1-4e55-81e4-b572ae7e3e63'), 3, 'Nibbles', 1, 10, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('4f2e0e7c-7f72-4429-8bc3-5f78c3e0d1c7'), 3, 'Draw', 1, 11, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
    cockSpace:auto_increment { uuid.fromstr('da28e05a-f033-4d90-8d6c-fdc14d4b9a38'), 3, 'Cancel', 1, 12, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_fights_schedules_space()
  local fightScheduleSpace = box.space.fights_schedules_space
  local seq_fight_schedule_id = box.sequence.seq_fight_schedule_id

  if fightScheduleSpace == nil then
    if seq_fight_schedule_id ~= nil then
      seq_fight_schedule_id:drop()
    end

    box.schema.sequence.create('seq_fight_schedule_id', { start = 1 })

    local format = {
      { "id",                  "unsigned" },
      { "fight_schedule_uuid", "uuid",     is_nullable = true },

      { "channel_id",          "number",   is_nullable = false, default = 1 },
      { "schedule_name",       "string",   is_nullable = false },
      { "schedule_desc",       "string",   is_nullable = false },
      { "schedule_date",       "datetime", is_nullable = false },
      { "status_id",           "number",   is_nullable = false, default = 1 },
      { "order",               "number",   is_nullable = false, default = 1 },

      { "created_by",          "number",   is_nullable = false },
      { "created_at",          "datetime", is_nullable = false },
      { "updated_by",          "number",   is_nullable = true },
      { "updated_at",          "datetime", is_nullable = true },
      { "deleted_by",          "number",   is_nullable = true },
      { "deleted_at",          "datetime", is_nullable = true },
    }

    fightScheduleSpace = box.schema.create_space('fights_schedules_space', { format = format, id = 1014 })

    fightScheduleSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_fight_schedule_id',
      if_not_exists = true
    })

    fightScheduleSpace:create_index('fight_schedule_uuid', {
      parts = { { 'fight_schedule_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })
  end
end

local function init_fights_schedules_details_space()
  local fightScheduleDetailSpace = box.space.fights_schedules_details_space
  local seq_fight_schedule_detail_id = box.sequence.seq_fight_schedule_detail_id

  if fightScheduleDetailSpace == nil then
    if seq_fight_schedule_detail_id ~= nil then
      seq_fight_schedule_detail_id:drop()
    end

    box.schema.sequence.create('seq_fight_schedule_detail_id', { start = 1 })

    local format = {
      { "id",                         "unsigned" },
      { "fight_schedule_detail_uuid", "uuid",     is_nullable = true },

      { "fight_schedule_id",          "number",   is_nullable = false, default = 1 },
      -- The idea check this https://www.facebook.com/photo/?fbid=363027626101404&set=pcb.363028292768004
      { "fight_no",                   "number",   is_nullable = false, default = 1 }, -- Match number ordering
      { "red_cock_id",                "number",   is_nullable = false },              -- Meron corner id
      { "red_corner_name",            "string",   is_nullable = false },              --Meron cock name
      { "red_corner_weight",          "number",   is_nullable = false },              --Meron weight
      { "blue_cock_id",               "number",   is_nullable = false },              -- Wala corner id
      { "blue_corner_name",           "string",   is_nullable = false },              -- Wala corner name
      { "blue_corner_weight",         "number",   is_nullable = false },
      { "schedule_date",              "datetime", is_nullable = false },
      { "status_id",                  "number",   is_nullable = false, default = 1 },
      { "order",                      "number",   is_nullable = false, default = 1 },

      { "created_by",                 "number",   is_nullable = false },
      { "created_at",                 "datetime", is_nullable = false },
      { "updated_by",                 "number",   is_nullable = true },
      { "updated_at",                 "datetime", is_nullable = true },
      { "deleted_by",                 "number",   is_nullable = true },
      { "deleted_at",                 "datetime", is_nullable = true },
    }

    fightScheduleDetailSpace = box.schema.create_space('fights_schedules_details_space', { format = format, id = 1015 })

    fightScheduleDetailSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_fight_schedule_detail_id',
      if_not_exists = true
    })

    fightScheduleDetailSpace:create_index('fight_schedule_detail_uuid', {
      parts = { { 'fight_schedule_detail_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })
  end
end

local function init_fights_odds_space()
  local fightOddSpace = box.space.fights_odds_space
  local seq_fight_odd_id = box.sequence.seq_fight_odd_id

  if fightOddSpace == nil then
    if seq_fight_odd_id ~= nil then
      seq_fight_odd_id:drop()
    end

    box.schema.sequence.create('seq_fight_odd_id', { start = 1 })

    local format = {
      { "id",             "unsigned" },
      { "fight_odd_uuid", "uuid",     is_nullable = true },

      { "channel_id",     "number",   is_nullable = false, default = 1 },
      { "round_id",       "number",   is_nullable = false },
      { "bet_id",         "number",   is_nullable = false },
      { "red_cock_id",    "number",   is_nullable = false }, -- E.g. Meron
      { "red_odd",        "decimal",  is_nullable = false }, -- E.g. 2.10
      { "draw_odd",       "decimal",  is_nullable = false }, -- E.g. draw 0.8
      { "blue_cock_id",   "number",   is_nullable = false }, -- E.g. Wala
      { "blue_odd",       "decimal",  is_nullable = false }, -- E.g. 1.95
      { "status_id",      "number",   is_nullable = false, default = 1 },
      { "order",          "number",   is_nullable = false, default = 1 },

      { "created_by",     "number",   is_nullable = false },
      { "created_at",     "datetime", is_nullable = false },
      { "updated_by",     "number",   is_nullable = true },
      { "updated_at",     "datetime", is_nullable = true },
      { "deleted_by",     "number",   is_nullable = true },
      { "deleted_at",     "datetime", is_nullable = true },
    }

    fightOddSpace = box.schema.create_space('fights_odds_space', { format = format, id = 1016 })

    fightOddSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_fight_odd_id',
      if_not_exists = true
    })

    fightOddSpace:create_index('fight_odd_uuid', {
      parts = { { 'fight_odd_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })

    fightOddSpace:create_index('round_id', {
      parts = { { 'round_id', 'number' } },
      unique = false,
      if_not_exists = true
    })

    fightOddSpace:create_index('bet_id', {
      parts = { { 'bet_id', 'number' } },
      unique = false,
      if_not_exists = true
    })

  end

  if fightOddSpace:len() == 0 then
    fightOddSpace:auto_increment {
      uuid.fromstr('2f3e4d5c-6b7a-4890-9f8e-7d6c5b4a3f2e'),
      1, 1, 1, 1,
      decimal.new(2.1),
      decimal.new(8),
      2,
      decimal.new(1.9),
      1, 1, 1,
      datetime.parse('2024-08-30T09:55:30.301612345Z'),
      box.NULL, box.NULL, box.NULL, box.NULL
    }
  end
end

local function init_bets_types_space()
  local betsTypesSpace = box.space.bets_types_space
  local seq_bet_type_id = box.sequence.seq_bet_type_id

  if betsTypesSpace == nil then
    if seq_bet_type_id ~= nil then
      seq_bet_type_id:drop()
    end

    box.schema.sequence.create('seq_bet_type_id', { start = 1 })

    local format = {
      { "id",            "unsigned" },
      { "bet_type_uuid", "uuid",     is_nullable = true },

      { "bet_type_name", "string",   is_nullable = false },
      { "channel_id",    "number",   is_nullable = false, default = 1 },
      { "cock_id",       "number",   is_nullable = false, default = 1 },
      { "status_id",     "number",   is_nullable = false, default = 1 },
      { "order",         "number",   is_nullable = false, default = 1 },

      { "created_by",    "number",   is_nullable = false },
      { "created_at",    "datetime", is_nullable = false },
      { "updated_by",    "number",   is_nullable = true },
      { "updated_at",    "datetime", is_nullable = true },
      { "deleted_by",    "number",   is_nullable = true },
      { "deleted_at",    "datetime", is_nullable = true },
    }

    betsTypesSpace = box.schema.create_space('bets_types_space', { format = format, id = 1017 })

    betsTypesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_bet_type_id',
      if_not_exists = true
    })

    betsTypesSpace:create_index('bet_type_uuid', {
      parts = { { 'bet_type_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })

    betsTypesSpace:create_index('channel_id',
      {
        parts = { { 'channel_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

  end

  if betsTypesSpace:len() == 0 then
      -- Channel 1 (Main characters)
      betsTypesSpace:auto_increment { uuid.fromstr('3d2c1b0a-9f8e-4d7c-6b5a-4f3e2d1c0b9a'), 'Tom Corner', 1, 1, 1, 1, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('7e6d5c4b-3a2f-4190-8e7d-6c5b4a3f2e1d'), 'Jerry Corner', 1, 2, 1, 2, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('f1e0d9c8-b7a6-4958-7f6e-5d4c3b2a1f0e'), 'Draw Corner', 1, 3, 1, 3, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('beeca01d-0ec6-4078-bb60-c328bf9b2a92'), 'Cancel Corner', 1, 4, 1, 4, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }

      -- Channel 2 (Supporting characters)
      betsTypesSpace:auto_increment { uuid.fromstr('d2a1c2e3-6f2b-4f8d-bd50-e5d0bca86bcf'), 'Spike Corner', 2, 5, 1, 1, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('e4f8d5c6-3b4e-4c9d-bd12-b2d21d6e9f3e'), 'Toodles Corner', 2, 6, 1, 2, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('b4e8f9c5-7e2b-4e3d-ae13-193e2f91c88c'), 'Draw Corner', 2, 7, 1, 3, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('c6b2f3a4-5d1e-4b3f-bf67-f4e24c9a0e4f'), 'Cancel Corner', 2, 8, 1, 4, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }

      -- Channel 3 (Famous characters / recurring guests)
      betsTypesSpace:auto_increment { uuid.fromstr('c3a1e4d2-9f8e-4d7c-8a90-b3d2d3e6bff5'), 'Butch Corner', 3, 9, 1, 1, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('b5e7d9a3-4e2f-4f8d-8b65-8f56d3d2991e'), 'Nibbles Corner', 3, 10, 1, 2, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('f2e3d4a5-1e1c-4a4d-a234-93eaf3d2512a'), 'Draw Corner', 3, 11, 1, 3, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
      betsTypesSpace:auto_increment { uuid.fromstr('d1f2e3c4-2d3a-4b4d-bb76-7f8c4b7f1212'), 'Cancel Corner', 3, 12, 1, 4, 1, datetime.parse('2024-08-30T09:55:30.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_defaults_ratelimits_space() --[[ for general ratelimits on players --]]
  local defaults_ratelimits_space = box.space.defaults_ratelimits_space
  local seq_default_ratelimit_id = box.sequence.seq_default_ratelimit_id

  if defaults_ratelimits_space == nil then
    if seq_default_ratelimit_id ~= nil then
      seq_default_ratelimit_id:drop()
    end

    box.schema.sequence.create('seq_default_ratelimit_id', { start = 1 })

    local format = {
      { "id",                     "unsigned" },
      { "default_ratelimit_uuid", "uuid",     is_nullable = true },

      { "bet_type_id",            "number",   is_nullable = false },
      { "min_amount",             "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "warning_amount",         "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "max_amount",             "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "min_amount_each_bet",    "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "max_amount_each_bet",    "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "currency_id",            "number",   is_nullable = false },

      { "status_id",              "number",   is_nullable = false, default = 1 },
      { "order",                  "number",   is_nullable = false, default = 1 },
      { "created_by",             "number",   is_nullable = false },
      { "created_at",             "datetime", is_nullable = false },
      { "updated_by",             "number",   is_nullable = true },
      { "updated_at",             "datetime", is_nullable = true },
      { "deleted_by",             "number",   is_nullable = true },
      { "deleted_at",             "datetime", is_nullable = true },
    }

    defaults_ratelimits_space = box.schema.create_space('defaults_ratelimits_space', { format = format, id = 1018 })

    defaults_ratelimits_space:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_default_ratelimit_id',
      if_not_exists = true
    })
    defaults_ratelimits_space:create_index('currency_id',
      {
        parts = { { 'currency_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
  end

  if defaults_ratelimits_space:len() == 0 then
    defaults_ratelimits_space:auto_increment { uuid.fromstr('cdc4183d-b98a-48af-a97b-c8e26a9c7dde'), 1, decimal.new(1000), decimal.new(700000), decimal.new(1000000), decimal.new(1000),decimal.new(500000), 1, 1, 1, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('9c0672fd-f9d1-4d8c-a38e-b3f02dfc4494'), 2, decimal.new(1000), decimal.new(700000), decimal.new(1000000), decimal.new(1000),decimal.new(500000), 1, 1, 2, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('9fbf539a-ee13-4f87-b521-651abfeb9885'), 3, decimal.new(1000), decimal.new(700000), decimal.new(1000000), decimal.new(1000),decimal.new(500000), 1, 1, 3, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    -- USD
    defaults_ratelimits_space:auto_increment { uuid.fromstr('0fae8caa-1140-4b8f-9461-e9beccd1b5d8'), 4, decimal.new(0.25), decimal.new(700), decimal.new(1000), decimal.new(0.25), decimal.new(500), 2, 1, 4, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('5cf65188-9754-4b44-b59a-6ae175b41a72'), 5, decimal.new(0.25), decimal.new(700), decimal.new(1000), decimal.new(0.25), decimal.new(500), 2, 1, 5, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('9d6ed0db-2181-4820-b379-1fe2190bb5b3'), 6, decimal.new(0.25), decimal.new(700), decimal.new(1000), decimal.new(0.25), decimal.new(500), 2, 1, 6, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }

    -- VND
    defaults_ratelimits_space:auto_increment { uuid.fromstr('632d6ee4-8074-4068-b99d-23b031f8848c'), 7, decimal.new(10000), decimal.new(9000000), decimal.new(100000000), decimal.new(10000), decimal.new(50000000), 3, 1, 7, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('50d27676-92bc-4abf-90fb-55ed22068ed8'), 8, decimal.new(10000), decimal.new(9000000), decimal.new(100000000), decimal.new(10000), decimal.new(50000000), 3, 1, 8, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('c7d01962-58f7-4ba1-9165-35dc488bb589'), 9, decimal.new(10000), decimal.new(9000000), decimal.new(100000000), decimal.new(10000), decimal.new(50000000), 3, 1, 9, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    -- THB
    defaults_ratelimits_space:auto_increment { uuid.fromstr('8da11924-a030-4c05-95c6-9f7d5990d4e2'), 10, decimal.new(10), decimal.new(9000), decimal.new(10000), decimal.new(10), decimal.new(5000), 4, 1, 10, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('5cdc5507-41d1-48b6-ab47-14f73e468146'), 11, decimal.new(10), decimal.new(9000), decimal.new(10000), decimal.new(10), decimal.new(5000), 4, 1, 11, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('a2b12631-b2f1-40e7-ac5a-9d2cd43e6211'), 12, decimal.new(10), decimal.new(9000), decimal.new(10000), decimal.new(10), decimal.new(5000), 4, 1, 12, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    -- CNY
    defaults_ratelimits_space:auto_increment { uuid.fromstr('c7956e7d-f3f9-4147-9efa-e72785c8d9d2'), 13, decimal.new(2), decimal.new(9000), decimal.new(10000), decimal.new(2), decimal.new(5000), 5, 1, 13, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('cebcaa7e-b18e-4afe-94d8-ea3cde35fb54'), 14, decimal.new(2), decimal.new(9000), decimal.new(10000), decimal.new(2), decimal.new(5000), 5, 1, 14, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
    defaults_ratelimits_space:auto_increment { uuid.fromstr('6096e99d-b0e0-4f24-9f8d-12ada74d70df'), 15, decimal.new(2), decimal.new(9000), decimal.new(10000), decimal.new(2), decimal.new(5000), 5, 1, 15, 1, datetime.parse('2024-08-29T16:55:10.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_member_spaces()
  local memberSpace = box.space.members_space
  local seq_member_id = box.sequence.seq_member_id

  if memberSpace == nil then
    -- Drop the existing sequence if it exists
    if seq_member_id ~= nil then
      seq_member_id:drop()
    end

    -- Create a new sequence for player IDs
    box.schema.sequence.create('seq_member_id', { start = 1 })

    -- Define the space format
    local format = {
      { "id",               "unsigned" },
      { "member_uuid",      "uuid",     is_nullable = false },
      { "first_name",       "string",   is_nullable = false },
      { "last_name",        "string",   is_nullable = false },
      { "user_name",        "string",   is_nullable = false },
      { "password",         "string",   is_nullable = false },
      { "email",            "string",   is_nullable = false },
      { "login_session",    "string",   is_nullable = true },
      { "profile_photo",    "string",   is_nullable = true },
      { "member_alias",     "string",   is_nullable = true },
      { "phone_number",     "string",   is_nullable = true },
      { "role_id",          "number",   is_nullable = true, default = 1 }, -- Default to 1
      { "member_avatar_id", "number",   is_nullable = true },
      { "commission",       "decimal",  is_nullable = true, default = decimal.new(0.00) },
      { "token_version",    "number",   is_nullable = true, default = 1 },
      { "status_id",        "number",   is_nullable = true, default = 1 },
      { "order",            "number",   is_nullable = true, default = 1 },
      { "created_by",       "number",   is_nullable = false },
      { "created_at",       "datetime", is_nullable = false }, -- Must not be omitted
      { "updated_by",       "number",   is_nullable = true },
      { "updated_at",       "datetime", is_nullable = true },
      { "deleted_by",       "number",   is_nullable = true },
      { "deleted_at",       "datetime", is_nullable = true },
    }

    -- Create the space with the specified format
    memberSpace = box.schema.create_space('members_space', { format = format, id = 1019 })

    -- Create indexes
    memberSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_member_id',
      if_not_exists = true,
    })
    memberSpace:create_index('member_uuid', {
      parts = { { 'member_uuid', 'uuid' } },
      if_not_exists = true,
      unique = true,
    })
    memberSpace:create_index('member_alias', {
      parts = { { 'member_alias', 'string' } },
      if_not_exists = true,
      unique = true,
    })

    -- Set sequence starting value
    box.sequence.seq_member_id:set(0)
  end

  -- Insert a default record if the space is empty
local memberSpace = box.space.members_space

if memberSpace:len() == 0 then
    -- PLAYER 01
    memberSpace:auto_increment {
      uuid.fromstr('e2234678-c0c4-4d0b-9179-5cb5a2ece01f'),
      'Player',
      '01',
      'PLAYER001',
      '123456',
      'player.001@gmail.com',
      'bdeb581454a4441784be1e355faeab63',
      'player001.png',
      'KM001',
      '010123123',
      1,                    -- role_player_id
      nil,                  -- player_avatar_id
      decimal.new(0.00),    -- commission
      1,                    -- token_version
      1,                    -- status_id
      1,                    -- order
      1,                    -- created_by
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil,
    }

    -- PLAYER 02
    memberSpace:auto_increment {
      uuid.fromstr('7cbbc6fe-9e4d-4d41-b7f7-46a09c3f9a02'),
      'Player',
      '02',
      'PLAYER002',
      '123456',
      'player.002@gmail.com',
      nil,
      'player002.png',
      'KM002',
      '010222222',
      1,
      nil,
      decimal.new(0.00),
      1,
      1,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil,
    }

    -- PLAYER 03
    memberSpace:auto_increment {
      uuid.fromstr('eaa77030-63e8-46f3-b9b0-97f2debd2c03'),
      'Player',
      '03',
      'PLAYER003',
      '123456',
      'player.003@gmail.com',
      nil,
      'player003.png',
      'KM003',
      '010333333',
      1,
      nil,
      decimal.new(0.00),
      1,
      1,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil,
    }

    -- PLAYER 04
    memberSpace:auto_increment {
      uuid.fromstr('4ef0092b-2b48-4ff7-9d06-2cbef3c3f104'),
      'Player',
      '04',
      'PLAYER004',
      '123456',
      'player.004@gmail.com',
      nil,
      'player004.png',
      'KM004',
      '010444444',
      1,
      nil,
      decimal.new(0.00),
      1,
      1,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil,
    }

  end

end

local function init_coins_assets_space()
  local assetsSpace = box.space.coins_assets_space
  local seq_coin_asset_id = box.sequence.seq_coin_asset_id

  if assetsSpace == nil then
    if seq_coin_asset_id ~= nil then
      seq_coin_asset_id:drop()
    end

    box.schema.sequence.create('seq_coin_asset_id', { start = 1 })

    local format = {
      { "id",              "unsigned" },
      { "coin_asset_uuid", "uuid",     is_nullable = true },
      { "file_name",       "string",   is_nullable = true },
      { "css_cls",         "string",   is_nullable = true },
      { "img_number",      "number",   is_nullable = true },
      { "status",          "boolean",  is_nullable = true, default = true },

      { "created_by",      "number",   is_nullable = false },
      { "created_at",      "datetime", is_nullable = false },
      { "updated_by",      "number",   is_nullable = true },
      { "updated_at",      "datetime", is_nullable = true },
      { "deleted_by",      "number",   is_nullable = true },
      { "deleted_at",      "datetime", is_nullable = true },

    }

    assetsSpace = box.schema.create_space('coins_assets_space', { format = format, id = 1020 })

    assetsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_coin_asset_id',
      if_not_exists = true
    })

    assetsSpace:create_index('coin_asset_uuid',
      {
        parts = { { 'coin_asset_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if assetsSpace:len() == 0 then
    -- Insert default records
    assetsSpace:auto_increment { uuid.fromstr('f4996c3a-6646-4184-80be-3e60ac100fe9'), 'coin1.svg', 'coin1', 1, true, 1, datetime.parse('2024-06-19T17:32:40.301612345Z'), nil, nil, nil, nil }
    assetsSpace:auto_increment { uuid.fromstr('def022c0-9805-4545-98f5-30f380ffb88e'), 'coin2.svg', 'coin2', 2, true, 1, datetime.parse('2024-06-19T17:32:41.301612345Z'), nil, nil, nil, nil }
    assetsSpace:auto_increment { uuid.fromstr('89719801-4fe9-4b05-992f-216962894719'), 'coin3.svg', 'coin3', 3, true, 1, datetime.parse('2024-06-19T17:33:42.301612345Z'), nil, nil, nil, nil }
    assetsSpace:auto_increment { uuid.fromstr('d4c789f6-4d9d-4b35-9dec-eaf374889fc9'), 'coin4.svg', 'coin4', 4, true, 1, datetime.parse('2024-06-19T17:33:43.301612345Z'), nil, nil, nil, nil }
    assetsSpace:auto_increment { uuid.fromstr('3c118a12-68d5-45b2-944e-daf47ecdd4b6'), 'coin5.svg', 'coin5', 5, true, 1, datetime.parse('2024-06-19T17:33:00.301612345Z'), nil, nil, nil, nil }
    assetsSpace:auto_increment { uuid.fromstr('211cea46-9bfe-450f-bed7-5199d64b1911'), 'coin6.svg', 'coin6', 6, true, 1, datetime.parse('2024-06-19T17:34:10.301612345Z'), nil, nil, nil, nil }
    assetsSpace:auto_increment { uuid.fromstr('6107a582-add7-418d-ab06-972b6d160066'), 'coin7.svg', 'coin7', 7, true, 1, datetime.parse('2024-06-19T17:35:00.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_coins_space()
  local coinsSpace = box.space.coins_space
  local seq_coin_id = box.sequence.seq_coin_id

  if coinsSpace == nil then
    if seq_coin_id ~= nil then
      seq_coin_id:drop()
    end

    box.schema.sequence.create('seq_coin_id', { start = 1 })

    local format = {
      { "id",            "unsigned" },
      { "coin_uuid",     "uuid",     is_nullable = false },
      { "group_id",      "number",   is_nullable = true }, -- ជំពាក់ space​ `coins_groups_space` អាចអត់បាន
      { "currency_id",   "number",   is_nullable = false },
      { "label",         "string",   is_nullable = false },
      { "value",         "number",   is_nullable = false },
      { "color_code",    "string",   is_nullable = true },

      { "coin_asset_id", "number",   is_nullable = true }, -- from coin_assets_space, field id
      { "is_active",     "boolean",  is_nullable = false, default = true },
      { "status_id",     "number",   is_nullable = true,  default = 1 },
      { "order",         "number",   is_nullable = true,  default = 1 },

      { "created_by",    "number",   is_nullable = true },
      { "created_at",    "datetime", is_nullable = true },
      { "updated_at",    "datetime", is_nullable = true },
      { "updated_by",    "number",   is_nullable = true },
      { "deleted_at",    "datetime", is_nullable = true },
      { "deleted_by",    "number",   is_nullable = true }

    }

    coinsSpace = box.schema.create_space('coins_space', { format = format, id = 1021 })

    coinsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_coin_id',
      if_not_exists = true
    })
    coinsSpace:create_index('coin_uuid_currency_id_coin_asset_id',
      {
        parts = { { 'coin_uuid', 'uuid' }, { 'currency_id', 'number' }, { 'coin_asset_id', 'number' } },
        if_not_exists = true,
        unique = true
      }
    )
    coinsSpace:create_index('coin_uuid',
      {
        parts = { { 'coin_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      })

    coinsSpace:create_index('currency_id',
      {
        parts = { { 'currency_id', 'number' } },
        if_not_exists = true,
        unique = false
      })
    coinsSpace:create_index('coin_asset_id',
      {
        parts = { { 'coin_asset_id', 'number' } },
        if_not_exists = true,
        unique = false
      })
  end

  if coinsSpace:len() == 0 then
    -- Insert default records
    -- KHR
    coinsSpace:auto_increment { uuid.fromstr('962c908e-70b8-4e5a-8786-d23de7b8ed6f'), 1, 1, '៥ពាន់', 5000, 'blue', 1, true, 1, 1, 1, datetime.parse('2024-05-28T17:20:30.301612345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('3783ff01-0732-4094-b2fc-e4b33854c80d'), 1, 1, '២ម៉ឺន', 20000, 'green', 2, true, 1, 2, 1, datetime.parse('2024-05-28T17:21:51.798111345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('99241ae8-cb4f-42c9-8890-320bdf98bd4e'), 1, 1, '៥ម៉ឺន', 50000, 'brown', 3, true, 1, 3, 1, datetime.parse('2024-05-28T17:22:45.716791831Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('71a3f6f0-7daf-450b-8354-a08853a1f63e'), 1, 1, '១០ម៉ឺន', 100000, 'yellow', 4, true, 1, 4, 1, datetime.parse('2024-05-28T17:23:25.692623639Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('8e489b33-2e45-41b6-b96f-1681c7273bd8'), 1, 1, '២០ម៉ឺន', 200000, 'dark-red', 5, true, 1, 5, 1, datetime.parse('2024-05-28T17:24:04.886882506Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('b7753b00-c3a0-4a37-85b5-57bc735d4a19'), 1, 1, '៤០ម៉ឺន', 400000, 'light-blue', 6, true, 1, 6, 1, datetime.parse('2024-05-28T17:24:40.479083523Z'), nil, nil, nil, nil }

    -- USD
    coinsSpace:auto_increment { uuid.fromstr('ae005c7a-0248-4f7c-874b-31cca485fc49'), 2, 2, '1$', 1, 'blue', 1, true, 1, 1, 1, datetime.parse('2024-02-04T17:20:30.301612345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('7be5f801-4b77-4dc6-b5f1-2622f378d896'), 2, 2, '5$', 5, 'green', 2, true, 1, 2, 1, datetime.parse('2024-02-04T17:21:51.798111345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('e60a110e-5f1d-4aab-bc07-31683e192e34'), 2, 2, '10$', 10, 'brown', 3, true, 1, 3, 1, datetime.parse('2024-02-04T17:22:45.716791831Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('b0c3e569-e7b7-49b9-b5ef-7cf51cfc8f7e'), 2, 2, '25$', 25, 'yellow', 4, true, 1, 4, 1, datetime.parse('2024-02-04T17:23:25.692623639Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('2e6866b4-2f55-4402-a5d3-a75720deab23'), 2, 2, '50$', 50, 'dark-red', 5, true, 1, 5, 1, datetime.parse('2024-02-04T17:24:04.886882506Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('d21d3d8d-b0ee-4f22-b455-7caf43c59829'), 2, 2, '100$', 100, 'light-blue', 6, true, 1, 6, 1, datetime.parse('2024-02-04T17:24:40.479083523Z'), nil, nil, nil, nil }

    -- VND
    coinsSpace:auto_increment { uuid.fromstr('9b1deb4d-3b7d-4bad-9bbd-2b594ac3816d'), 3, 3, '10d', 10000, 'blue', 1, true, 1, 1, 1, datetime.parse('2024-02-04T17:20:30.301612345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('c04a1fe0-9196-4138-9ee0-486f3f2e4482'), 3, 3, '100d', 100000, 'green', 2, true, 1, 2, 1, datetime.parse('2024-02-04T17:21:51.798111345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('3e463154-8834-4b64-81f7-d89cd696a8cf'), 3, 3, '500d', 500000, 'brown', 3, true, 1, 3, 1, datetime.parse('2024-02-04T17:22:45.716791831Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('b2e85b1a-9e2b-4ee3-9d39-bc2b0c9ef005'), 3, 3, '1000d', 1000000, 'yellow', 4, true, 1, 4, 1, datetime.parse('2024-02-04T17:23:25.692623639Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('6d7c55a8-9b43-4d24-8562-3d55b1499707'), 3, 3, '3000d', 3000000, 'dark-red', 5, true, 1, 5, 1, datetime.parse('2024-02-04T17:24:04.886882506Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('f18e6f4b-c273-4e69-9af8-72e10d22371d'), 3, 3, '7000d', 7000000, 'light-blue', 6, true, 1, 6, 1, datetime.parse('2024-02-04T17:24:40.479083523Z'), nil, nil, nil, nil }

    -- THB
    coinsSpace:auto_increment { uuid.fromstr('a39732e8-8f49-4f4f-8be6-92d0865f9f96'), 4, 4, '10฿', 10, 'blue', 1, true, 1, 1, 1, datetime.parse('2024-02-04T17:20:30.301612345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('0e2f45dd-f425-4f2a-99dc-ba6aa3b7e727'), 4, 4, '50฿', 50, 'green', 2, true, 1, 2, 1, datetime.parse('2024-02-04T17:21:51.798111345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('8e0a9b7b-6b3f-4f6d-9c4b-39d5dd62a025'), 4, 4, '100฿', 100, 'brown', 3, true, 1, 3, 1, datetime.parse('2024-02-04T17:22:45.716791831Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('4e93c0ef-a7ec-4819-8664-dfb055889212'), 4, 4, '500฿', 500, 'yellow', 4, true, 1, 4, 1, datetime.parse('2024-02-04T17:23:25.692623639Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('0edd6d2a-d5a5-4d86-9326-5cdf29b25764'), 4, 4, '2000฿', 2000, 'dark-red', 5, true, 1, 5, 1, datetime.parse('2024-02-04T17:24:04.886882506Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('76e25e53-2ef1-4d83-8918-04b69b143b12'), 4, 4, '5000฿', 5000, 'light-blue', 6, true, 1, 6, 1, datetime.parse('2024-02-04T17:24:40.479083523Z'), nil, nil, nil, nil }

    -- CNY
    coinsSpace:auto_increment { uuid.fromstr('2f4a9025-8519-4928-9916-7a577eb25c03'), 5, 5, '2¥', 2, 'blue', 1, true, 1, 1, 1, datetime.parse('2024-02-04T17:20:30.301612345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('b7cc7e1b-0b9b-4c8a-9436-d93464201e8a'), 5, 5, '10¥', 10, 'green', 2, true, 1, 2, 1, datetime.parse('2024-02-04T17:21:51.798111345Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('1d6f9664-9e7d-4a7d-86e3-6d2f7c60cd2c'), 5, 5, '20¥', 20, 'brown', 3, true, 1, 2, 1, datetime.parse('2024-02-04T17:22:45.716791831Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('6cbc23f5-0f99-4f6f-9f61-3a9e7d705f9e'), 5, 5, '100¥', 100, 'yellow', 4, true, 1, 4, 1, datetime.parse('2024-02-04T17:23:25.692623639Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('ea9eaac9-d289-41d9-80c9-2501c5afd782'), 5, 5, '400¥', 400, 'dark-red', 5, true, 1, 5, 1, datetime.parse('2024-02-04T17:24:04.886882506Z'), nil, nil, nil, nil }
    coinsSpace:auto_increment { uuid.fromstr('5f7fd8c6-3f8c-4c9a-80b6-9e7b5d09c9f8'), 5, 5, '1000¥', 1000, 'light-blue', 6, true, 1, 6, 1, datetime.parse('2024-02-04T17:24:40.479083523Z'), nil, nil, nil, nil }
  end
end

local function init_members_balances_space()
  local membersBalancesSpace = box.space.members_balances_space
  local seq_member_balance_id = box.sequence.seq_member_balance_id

  if membersBalancesSpace == nil then
    if seq_member_balance_id ~= nil then
      seq_member_balance_id:drop()
    end

    box.schema.sequence.create('seq_member_balance_id', { start = 1 })

    local format = {
      { "id",                  "unsigned" },
      { "member_balance_uuid", "uuid",     is_nullable = true },

      { "member_id",           "number",   is_nullable = false },
      { "currency_id",         "number",   is_nullable = false, default = 1 },
      { "balance",             "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "status_id",           "number",   is_nullable = false, default = 1 },
      { "order",               "number",   is_nullable = false, default = 1 },

      { "created_by",          "number",   is_nullable = false },
      { "created_at",          "datetime", is_nullable = false },
      { "updated_by",          "number",   is_nullable = true },
      { "updated_at",          "datetime", is_nullable = true },
      { "deleted_by",          "number",   is_nullable = true },
      { "deleted_at",          "datetime", is_nullable = true },
    }

    membersBalancesSpace = box.schema.create_space('members_balances_space', { format = format, id = 1022 })

    membersBalancesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_member_balance_id',
      if_not_exists = true
    })


    membersBalancesSpace:create_index('member_id',
      {
        parts = {
          { 'member_id', 'number' }
        },
        if_not_exists = true,
        unique = false
      }
    )
    box.sequence.seq_member_balance_id:set(0)
        -- Insert default balances if empty
    if membersBalancesSpace:len() == 0 then
      membersBalancesSpace:auto_increment {
        uuid.new(),      -- member_balance_uuid
        1,               -- member_id (PLAYER001)
        1,               -- currency_id (default)
        decimal.new(1000.0000), -- balance
        1,               -- status_id
        1,               -- order
        1,               -- created_by
        datetime.new(os.date('*t')), -- created_at
        nil,             -- updated_by
        nil,             -- updated_at
        nil,             -- deleted_by
        nil              -- deleted_at
      }

      membersBalancesSpace:auto_increment {
        uuid.new(),      -- member_balance_uuid
        2,               -- member_id (PLAYER002)
        1,
        decimal.new(2500.0000),
        1,
        1,
        1,
        datetime.new(os.date('*t')),
        nil,
        nil,
        nil,
        nil
      }

      membersBalancesSpace:auto_increment {
        uuid.new(),      -- member_balance_uuid
        3,               -- member_id (PLAYER003)
        1,
        decimal.new(5000.0000),
        1,
        1,
        1,
        datetime.new(os.date('*t')),
        nil,
        nil,
        nil,
        nil
      }
    end
  end
end

local function init_players_roles_space()
  local UsersRolesSpace = box.space.players_roles_space
  local seq_player_role_id = box.sequence.seq_player_role_id

  if UsersRolesSpace == nil then
    if seq_player_role_id ~= nil then
      seq_player_role_id:drop()
    end

    -- Create sequence for player roles
    box.schema.sequence.create('seq_player_role_id', { start = 1 })

    local format = {
      { "id",               "unsigned" },
      { "player_role_uuid", "uuid",     is_nullable = false },
      { "player_role_name", "string",   is_nullable = false },
      { "player_role_desc", "string",   is_nullable = false },
      { "status",           "boolean",  is_nullable = false },
      { "order",            "number",   is_nullable = true, default = 1 },
      { "created_by",       "number",   is_nullable = false },
      { "created_at",       "datetime", is_nullable = false },
      { "updated_by",       "number",   is_nullable = true },
      { "updated_at",       "datetime", is_nullable = true },
      { "deleted_by",       "number",   is_nullable = true },
      { "deleted_at",       "datetime", is_nullable = true },
    }

    -- Create the space
    UsersRolesSpace = box.schema.create_space('players_roles_space', { format = format, id = 1023 })

    -- Create indexes
    UsersRolesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_player_role_id',
      if_not_exists = true,
    })

    UsersRolesSpace:create_index('player_role_uuid', {
      parts = { { 'player_role_uuid', 'uuid' } },
      if_not_exists = true,
      unique = true
    })

    -- Reset the sequence to start from 0
    box.sequence.seq_player_role_id:set(0)
  end

  if UsersRolesSpace:len() == 0 then
    -- Insert records
    UsersRolesSpace:auto_increment {
      uuid.fromstr('9a6f17b3-f2d1-4df4-8ade-d1c8fbebdb97'),
      'admin',
      'Role Admin',
      true,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }

    UsersRolesSpace:auto_increment {
      uuid.fromstr('01918ff3-57fd-7c09-93a3-7077087550b0'),
      'moderator',
      'Role Moderator',
      true,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }

    UsersRolesSpace:auto_increment {
      uuid.fromstr('01918fdb-0d16-7b77-b0c6-dc681aada863'),
      'operator',
      'Role Operator',
      true,
      1,
      1,
      datetime.new(os.date('*t')),
      nil,
      nil,
      nil,
      nil
    }
  end
end

local function init_notifications_types_space()
  local notificationsTypesSpace = box.space.notifications_types_space
  local seq_notification_type_id = box.sequence.seq_notification_type_id

  if notificationsTypesSpace == nil then
    if seq_notification_type_id ~= nil then
      seq_notification_type_id:drop()
    end

    box.schema.sequence.create('seq_notification_type_id', { start = 1 })

    local format = {
      { "id",                     "unsigned" },
      { "notification_type_uuid", "uuid",     is_nullable = true },

      { "notification_type_name", "string",   is_nullable = false },
      { "remark",                 "string",   is_nullable = true },
      { "status_id",              "number",   is_nullable = true, default = 1 },
      { "order",                  "number",   is_nullable = true, default = 1 },

      { "created_by",             "number",   is_nullable = false },
      { "created_at",             "datetime", is_nullable = false },
      { "updated_by",             "number",   is_nullable = true },
      { "updated_at",             "datetime", is_nullable = true },
      { "deleted_by",             "number",   is_nullable = true },
      { "deleted_at",             "datetime", is_nullable = true },
    }

    notificationsTypesSpace = box.schema.create_space('notifications_types_space', { format = format, id = 1024

    })

    notificationsTypesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_notification_type_id',
      if_not_exists = true
    })
  end

  if notificationsTypesSpace:len() == 0 then
    notificationsTypesSpace:auto_increment { uuid.fromstr('e38214d6-a202-4c72-b820-e0b59bc0ffc9'), 'General', '', 1, 1, 1, datetime.parse('2024-06-22T17:10:20.301612345Z'), nil, nil, nil, nil }
    notificationsTypesSpace:auto_increment { uuid.fromstr('f99b9973-cf0c-4962-9572-8fbed37f291c'), 'Financial', '', 1, 2, 1, datetime.parse('2024-06-22T17:10:20.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_countries_space()
  local countriesSpace = box.space.countries_space
  local seq_country_id = box.sequence.seq_country_id

  if countriesSpace == nil then
    if seq_country_id ~= nil then
      seq_country_id:drop()
    end

    box.schema.sequence.create('seq_country_id', { start = 1 })

    local format = {
      { "id",           "unsigned" },
      { "country_uuid", "uuid",     is_nullable = true },

      { "language_id",  "number",   is_nullable = false },
      { "country_name", "string",   is_nullable = false },
      { "country_code", "number",   is_nullable = false },
      { "status_id",    "number",   is_nullable = true, default = 1 },
      { "order",        "number",   is_nullable = true, default = 1 },

      { "created_by",   "number",   is_nullable = false },
      { "created_at",   "datetime", is_nullable = false },
      { "updated_by",   "number",   is_nullable = true },
      { "updated_at",   "datetime", is_nullable = true },
      { "deleted_by",   "number",   is_nullable = true },
      { "deleted_at",   "datetime", is_nullable = true },
    }

    countriesSpace = box.schema.create_space('countries_space', { format = format, id = 1025 })

    countriesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_country_id',
      if_not_exists = true
    })
  end

  if countriesSpace:len() == 0 then
    countriesSpace:auto_increment { uuid.fromstr('e8fc66fd-f42f-4e1a-b6e5-254726acaf40'), 1, 'Cambodia', 855, 1, 1, 1, datetime.parse('2024-06-26T15:51:10.301612345Z'), nil, nil, nil, nil }
    countriesSpace:auto_increment { uuid.fromstr('77c1c92b-fd4b-482f-9dba-882f82f35bec'), 2, 'United Kingdom', 44, 1, 2, 1, datetime.parse('2024-06-26T15:51:20.301612345Z'), nil, nil, nil, nil }
    countriesSpace:auto_increment { uuid.fromstr('f3231590-fd7d-4176-8c0d-55b95efc2c2c'), 3, 'Thailand', 66, 1, 3, 1, datetime.parse('2024-06-26T15:51:30.301612345Z'), nil, nil, nil, nil }
    countriesSpace:auto_increment { uuid.fromstr('3d6012ec-0ceb-4094-aedc-5dc5d992908f'), 4, 'Vietnam', 84, 1, 4, 1, datetime.parse('2024-06-26T15:51:40.301612345Z'), nil, nil, nil, nil }
    countriesSpace:auto_increment { uuid.fromstr('3d6012ec-0ceb-4094-aedc-5dc5d992908f'), 5, 'China', 86, 1, 5, 1, datetime.parse('2024-06-26T15:51:40.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_players_avatars_space()
  local playersAvatarsSpace = box.space.players_avatars_space
  local seq_player_avatar_id = box.sequence.seq_player_avatar_id

  if playersAvatarsSpace == nil then
    if seq_player_avatar_id ~= nil then
      seq_player_avatar_id:drop()
    end

    box.schema.sequence.create('seq_player_avatar_id', { start = 1 })

    local format = {
      { "id",                 "unsigned" },
      { "player_avatar_uuid", "uuid",     is_nullable = true },

      { "path",               "string",   is_nullable = false },
      { "file_name",          "string",   is_nullable = false },
      { "status_id",          "number",   is_nullable = true, default = 1 },
      { "order",              "number",   is_nullable = true, default = 1 },

      { "created_by",         "number",   is_nullable = false },
      { "created_at",         "datetime", is_nullable = false },
      { "updated_by",         "number",   is_nullable = true },
      { "updated_at",         "datetime", is_nullable = true },
      { "deleted_by",         "number",   is_nullable = true },
      { "deleted_at",         "datetime", is_nullable = true },
    }

    playersAvatarsSpace = box.schema.create_space('players_avatars_space', { format = format, id = 1026 })

    playersAvatarsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_player_avatar_id',
      if_not_exists = true
    })
    playersAvatarsSpace:create_index('player_avatar_uuid',
      {
        parts = { { 'player_avatar_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if playersAvatarsSpace:len() == 0 then
    playersAvatarsSpace:auto_increment { uuid.fromstr('410898c4-6e04-4a8c-a4dd-8175ffd413fc'), '/images/avatars/', 'avatar_1.jpg', 1, 1, 1, datetime.parse('2024-06-28T11:13:10.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('1674be74-84eb-4df2-b60d-8e29328746a5'), '/images/avatars/', 'avatar_2.jpg', 1, 2, 1, datetime.parse('2024-06-28T11:13:20.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('9df69d1e-716d-4e46-bc49-98adb540de7c'), '/images/avatars/', 'avatar_3.jpg', 1, 3, 1, datetime.parse('2024-06-28T11:13:30.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('df7de320-6401-46ac-9010-128f0c2c63c5'), '/images/avatars/', 'avatar_4.jpg', 1, 4, 1, datetime.parse('2024-06-28T11:13:40.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('2a8198d9-a755-403d-a4b0-b19a95784efa'), '/images/avatars/', 'avatar_5.jpg', 1, 5, 1, datetime.parse('2024-06-28T11:13:50.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('c9641cd4-bc11-41c3-9f22-2e93cb742912'), '/images/avatars/', 'avatar_6.jpg', 1, 6, 1, datetime.parse('2024-06-28T11:14:00.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('15a3e5a3-8e2c-4b67-b6fa-6f4597f3475a'), '/images/avatars/', 'avatar_7.jpg', 1, 7, 1, datetime.parse('2024-06-28T11:14:10.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('3efbff60-9c62-4a60-9307-8697465a43db'), '/images/avatars/', 'avatar_8.jpg', 1, 8, 1, datetime.parse('2024-06-28T11:14:20.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('2ff32af2-e0db-4f1d-a2f7-5a9675175f66'), '/images/avatars/', 'avatar_9.jpg', 1, 9, 1, datetime.parse('2024-06-28T11:14:30.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('7b8a3940-8d7e-4eb6-bd02-03a7fb593c3d'), '/images/avatars/', 'avatar_10.jpg', 1, 10, 1, datetime.parse('2024-06-28T11:14:40.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('a2e23a14-1c2b-47ae-a712-6ab34a4df160'), '/images/avatars/', 'avatar_11.jpg', 1, 11, 1, datetime.parse('2024-06-28T11:14:50.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('40cdaed5-18ab-41bb-bf30-bb76db61eb91'), '/images/avatars/', 'avatar_12.jpg', 1, 12, 1, datetime.parse('2024-06-28T11:15:00.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('0d8e4760-6502-4910-b3ab-0a3e8d3e9af0'), '/images/avatars/', 'avatar_13.jpg', 1, 13, 1, datetime.parse('2024-06-28T11:15:10.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('00f174cb-0f70-474e-8166-d929b55f42c1'), '/images/avatars/', 'avatar_14.jpg', 1, 14, 1, datetime.parse('2024-06-28T11:15:20.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('1a948a74-1e54-490f-934b-18bca42e9d88'), '/images/avatars/', 'avatar_15.jpg', 1, 15, 1, datetime.parse('2024-06-28T11:15:30.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('04e6f803-8a2c-4202-94b4-965812e8b8cf'), '/images/avatars/', 'avatar_16.jpg', 1, 16, 1, datetime.parse('2024-06-28T11:15:40.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('b7412763-0e17-47b0-a469-47e0b97a0126'), '/images/avatars/', 'avatar_17.jpg', 1, 17, 1, datetime.parse('2024-06-28T11:15:50.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('9db4e442-9e95-4b32-9e7a-4d9cb9833077'), '/images/avatars/', 'avatar_18.jpg', 1, 18, 1, datetime.parse('2024-06-28T11:16:00.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('4029a29b-b5d8-4dd8-80d0-0644518466f4'), '/images/avatars/', 'avatar_19.jpg', 1, 19, 1, datetime.parse('2024-06-28T11:16:10.301612345Z'), nil, nil, nil, nil }
    playersAvatarsSpace:auto_increment { uuid.fromstr('f669fa58-3df2-4a73-a870-b44db20d0a31'), '/images/avatars/', 'avatar_20.jpg', 1, 20, 1, datetime.parse('2024-06-28T11:16:20.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_settings_music_space()
  local settingsMusicSpace = box.space.settings_music_space
  local seq_setting_music_id = box.sequence.seq_setting_music_id

  if settingsMusicSpace == nil then
    if seq_setting_music_id ~= nil then
      seq_setting_music_id:drop()
    end

    box.schema.sequence.create('seq_setting_music_id', { start = 1 })

    local format = {
      { "id",                 "unsigned" },
      { "setting_music_uuid", "uuid",     is_nullable = true },

      { "min_number",         "number",   is_nullable = true, default = 0 },
      { "default_number",     "number",   is_nullable = true, default = 50 },
      { "max_number",         "number",   is_nullable = true, default = 100 },
      { "status_id",          "number",   is_nullable = true, default = 1 },
      { "order",              "number",   is_nullable = true, default = 1 },

      { "created_by",         "number",   is_nullable = false },
      { "created_at",         "datetime", is_nullable = false },
      { "updated_by",         "number",   is_nullable = true },
      { "updated_at",         "datetime", is_nullable = true },
      { "deleted_by",         "number",   is_nullable = true },
      { "deleted_at",         "datetime", is_nullable = true },
    }

    settingsMusicSpace = box.schema.create_space('settings_music_space', { format = format, id = 1027 })

    settingsMusicSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_setting_music_id',
      if_not_exists = true
    })
    settingsMusicSpace:create_index('setting_music_uuid',
      {
        parts = { { 'setting_music_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if settingsMusicSpace:len() == 0 then
    settingsMusicSpace:auto_increment { uuid.fromstr('1ed69d51-0053-4351-8030-5f9bc7215dba'), 0, 50, 100, 1, 1, 1, datetime.parse('2024-06-29T14:19:10.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_settings_volume_space()
  local settingsVolumeSpace = box.space.settings_volume_space
  local seq_setting_volume_id = box.sequence.seq_setting_volume_id

  if settingsVolumeSpace == nil then
    if seq_setting_volume_id ~= nil then
      seq_setting_volume_id:drop()
    end

    box.schema.sequence.create('seq_setting_volume_id', { start = 1 })

    local format = {
      { "id",                  "unsigned" },
      { "setting_volume_uuid", "uuid",     is_nullable = true },

      { "min_number",          "number",   is_nullable = true, default = 0 },
      { "default_number",      "number",   is_nullable = true, default = 50 },
      { "max_number",          "number",   is_nullable = true, default = 100 },
      { "status_id",           "number",   is_nullable = true, default = 1 },
      { "order",               "number",   is_nullable = true, default = 1 },

      { "created_by",          "number",   is_nullable = false },
      { "created_at",          "datetime", is_nullable = false },
      { "updated_by",          "number",   is_nullable = true },
      { "updated_at",          "datetime", is_nullable = true },
      { "deleted_by",          "number",   is_nullable = true },
      { "deleted_at",          "datetime", is_nullable = true },
    }

    settingsVolumeSpace = box.schema.create_space('settings_volume_space', { format = format, id = 1028 })

    settingsVolumeSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_setting_volume_id',
      if_not_exists = true
    })
    settingsVolumeSpace:create_index('setting_volume_uuid',
      {
        parts = { { 'setting_volume_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if settingsVolumeSpace:len() == 0 then
    settingsVolumeSpace:auto_increment { uuid.fromstr('2b8d1a6a-de30-48e1-af2f-f4fcd530ad3d'), 0, 50, 100, 1, 1, 1, datetime.parse('2024-06-29T14:19:10.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_settings_voice_space()
  local settingsVoiceSpace = box.space.settings_voice_space
  local seq_setting_voice_id = box.sequence.seq_setting_voice_id

  if settingsVoiceSpace == nil then
    if seq_setting_voice_id ~= nil then
      seq_setting_voice_id:drop()
    end

    box.schema.sequence.create('seq_setting_voice_id', { start = 1 })

    local format = {
      { "id",                 "unsigned" },
      { "setting_voice_uuid", "uuid",     is_nullable = true },

      { "min_number",         "number",   is_nullable = true, default = 0 },
      { "default_number",     "number",   is_nullable = true, default = 50 },
      { "max_number",         "number",   is_nullable = true, default = 100 },
      { "status_id",          "number",   is_nullable = true, default = 1 },
      { "order",              "number",   is_nullable = true, default = 1 },

      { "created_by",         "number",   is_nullable = false },
      { "created_at",         "datetime", is_nullable = false },
      { "updated_by",         "number",   is_nullable = true },
      { "updated_at",         "datetime", is_nullable = true },
      { "deleted_by",         "number",   is_nullable = true },
      { "deleted_at",         "datetime", is_nullable = true },
    }

    settingsVoiceSpace = box.schema.create_space('settings_voice_space', { format = format, id = 1029 })
    settingsVoiceSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_setting_voice_id',
      if_not_exists = true
    })
    settingsVoiceSpace:create_index('setting_voice_uuid',
      {
        parts = { { 'setting_voice_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if settingsVoiceSpace:len() == 0 then
    settingsVoiceSpace:auto_increment { uuid.fromstr('758d9669-bd5f-496b-8323-478d5bf155b1'), 0, 50, 100, 1, 1, 1, datetime.parse('2024-06-29T14:19:10.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_players_ratelimits_space() --[[ for general ratelimits on players --]]
  local players_ratelimits_space = box.space.players_ratelimits_space
  local seq_player_ratelimit_id = box.sequence.seq_player_ratelimit_id

  if players_ratelimits_space == nil then
    if seq_player_ratelimit_id ~= nil then
      seq_player_ratelimit_id:drop()
    end

    box.schema.sequence.create('seq_player_ratelimit_id', { start = 1 })

    local format = {
      { "id",                    "unsigned" },
      { "player_ratelimit_uuid", "uuid",     is_nullable = true },

      { "player_id",             "number",   is_nullable = false },
      { "bet_type_id",           "number",   is_nullable = false },
      { "min_amount",            "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "warning_amount",        "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "max_amount",            "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "min_amount_each_bet",    "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "max_amount_each_bet",    "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "currency_id",           "number",   is_nullable = false },

      { "status_id",             "number",   is_nullable = false, default = 1 },
      { "order",                 "number",   is_nullable = false, default = 1 },
      { "created_by",            "number",   is_nullable = false },
      { "created_at",            "datetime", is_nullable = false },
      { "updated_by",            "number",   is_nullable = true },
      { "updated_at",            "datetime", is_nullable = true },
      { "deleted_by",            "number",   is_nullable = true },
      { "deleted_at",            "datetime", is_nullable = true },
    }

    players_ratelimits_space = box.schema.create_space('players_ratelimits_space', { format = format, id = 1030 })

    players_ratelimits_space:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_player_ratelimit_id',
      if_not_exists = true
    })
    players_ratelimits_space:create_index('currency_id',
      {
        parts = { { 'currency_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
  end
  if players_ratelimits_space:len() == 0 then
    players_ratelimits_space:auto_increment {
      uuid.fromstr('758d9669-bd5f-496b-8323-478d5bf155b1'),
      1, 1,
      decimal.new(0.0000),
      decimal.new(25000000.000),
      decimal.new(50000000.000),
      decimal.new(0.000),
      decimal.new(25000000.000),
      1, 1, 1,
      1,
      datetime.parse('2024-10-27T10:55:13Z'),
      nil, nil, nil, nil
    }
  end
end

local function init_players_audits_spaces()
  local playersAuditsSpace = box.space.players_audits_spaces
  local seq_user_audit_id = box.sequence.seq_player_audit_id

  if playersAuditsSpace == nil then
    if seq_user_audit_id ~= nil then
      seq_user_audit_id:drop()
    end

    box.schema.sequence.create('seq_player_audit_id', { start = 1 })

    local format = {
      { "id",                   "unsigned" },
      { "player_audit_uuid",    "uuid",     is_nullable = false },

      { "player_id",            "number",   is_nullable = false },
      { "player_audit_context", "string",   is_nullable = false },
      { "player_audit_desc",    "string",   is_nullable = false },
      { "audit_type_id",        "number",   is_nullable = false },
      { "user_agent",           "string",   is_nullable = false },
      { "operator",             "string",   is_nullable = false },
      { "ip",                   "string",   is_nullable = false },
      { "status_id",            "number",   is_nullable = false },
      { "order",                "number",   is_nullable = true, default = 1 },
      { "created_by",           "number",   is_nullable = false },
      { "created_at",           "datetime", is_nullable = false },
      { "updated_by",           "number",   is_nullable = true },
      { "updated_at",           "datetime", is_nullable = true },
      { "deleted_by",           "number",   is_nullable = true },
      { "deleted_at",           "datetime", is_nullable = true },
    }

    playersAuditsSpace = box.schema.create_space('players_audits_spaces', { format = format, id = 1031 })

    -- Create index
    playersAuditsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_player_audit_id',
      if_not_exists = true,
    })
    playersAuditsSpace:create_index('player_audit_uuid',
      {
        parts = { { 'player_audit_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      })
    box.sequence.seq_user_audit_id:set(0)
  end
end

local function init_tickets_space()
  local ticketsSpace = box.space.tickets_space
  local seq_ticket_id = box.sequence.seq_ticket_id

  if ticketsSpace == nil then
    if seq_ticket_id ~= nil then
      seq_ticket_id:drop()
    end

    box.schema.sequence.create('seq_ticket_id', { start = 1 })
    -- box.sequence.seq_ticket_id:next();

    local format = {
      { "id",                   "unsigned" },
      { "ticket_uuid",          "uuid",     is_nullable = true },
      { "ticket_no",            "string",   is_nullable = false }, -- E.g. TK001-AA003-BET120  any purpose for UI, many games industry format no base on top down membership nickname

      { "channel_id",           "number",   is_nullable = false }, --* The same TJ system can copy running many
      { "round_id",             "number",   is_nullable = false },
      { "player_id",            "number",   is_nullable = false },
      { "bet_type_id",          "number",   is_nullable = false },
      { "platform_id",          "number",   is_nullable = true }, --* private field save player device e.g. android, web
      { "domain_name",          "string",   is_nullable = true }, --* private field save domain player place betting
      { "total_amount",         "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "total_invalid_amount", "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "status_id",            "number",   is_nullable = false, default = 1 },
      { "order",                "number",   is_nullable = false, default = 1 },

      { "created_by",           "number",   is_nullable = false },
      { "created_at",           "datetime", is_nullable = false },
      { "updated_by",           "number",   is_nullable = true },
      { "updated_at",           "datetime", is_nullable = true },
      { "deleted_by",           "number",   is_nullable = true },
      { "deleted_at",           "datetime", is_nullable = true },
    }

    ticketsSpace = box.schema.create_space('tickets_space', { format = format, id = 1032 })

    ticketsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_ticket_id',
      if_not_exists = true
    })


    ticketsSpace:create_index('ticket_uuid',
      {
        parts = { { 'ticket_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
    ticketsSpace:create_index('player_id',
      {
        parts = { { 'player_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
    ticketsSpace:create_index('round_id',
      {
        parts = { { 'round_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

    box.sequence.seq_ticket_id:set(0)
  end
end

local function init_bets_space()
  local bets_space = box.space.bets_space
  local seq_bet_id = box.sequence.seq_bet_id

  if bets_space == nil then

    if seq_bet_id ~= nil then
      seq_bet_id:drop()
    end

  box.schema.sequence.create('seq_bet_id', { start = 1 })

  local format = {
    { "id",               "unsigned" },
    { "bet_uuid",         "uuid",     is_nullable = true },

    { "round_id",         "number",   is_nullable = false },
    { "table_id",          "number",   is_nullable = false },
    { "member_id",        "number",   is_nullable = false },

    { "currency_id",      "number",   is_nullable = false },
    { "bet_amount",       "decimal",  is_nullable = false, default = decimal.new(0) },
    { "payout_amount",    "decimal",  is_nullable = true,  default = decimal.new(0) },
    { "is_win",           "boolean",  is_nullable = false, default = false },

    -- Player cards & result
    { "member_cards",     "array",    is_nullable = false },
    { "member_score",     "number",   is_nullable = false },
    { "member_is_pok",    "boolean",  is_nullable = false },
    { "member_deng",      "number",   is_nullable = false },

    { "result_id",        "number",   is_nullable = true },

    { "status_id",        "number",   is_nullable = false, default = 1 },
    -- 1=bet, 2=settled

    { "created_by",       "number",   is_nullable = false },
    { "created_at",       "datetime", is_nullable = false },
    { "updated_by",       "number",   is_nullable = true },
    { "updated_at",       "datetime", is_nullable = true },
    { "deleted_by",       "number",   is_nullable = true },
    { "deleted_at",       "datetime", is_nullable = true },
  }

  bets_space = box.schema.create_space('bets_space', {
    id = 1033,
    format = format,
  })

  bets_space:create_index('id', {
    parts = { { 'id', 'unsigned' } },
    sequence = 'seq_bet_id',
    if_not_exists = true
  })

  bets_space:create_index('bet_uuid', {
    parts = { { 'bet_uuid', 'uuid' } },
    unique = true,
    if_not_exists = true
  })

  bets_space:create_index('round_id', {
    parts = { { 'round_id', 'number' } },
    unique = false,
  })

  bets_space:create_index('member_id', {
    parts = { { 'member_id', 'number' } },
     unique = false,
  })

  bets_space:create_index('table_id', {
    parts = { { 'table_id', 'number' } },
     unique = false,
  })

  bets_space:create_index('result_id', {
    parts = { { 'result_id', 'number' } },
     unique = false,
  })

  box.sequence.seq_bet_id:set(0)
  end
  
end


-- This spapce mirror/the same tickets_space the purpose aovoid slow
local function init_statements_space()
  local statements_space = box.space.statements_space
  local seq_statement_id = box.sequence.seq_statement_id


  if statements_space == nil then
    if seq_statement_id ~= nil then
      seq_statement_id:drop()
    end

    box.schema.sequence.create('seq_statement_id', { start = 1 })

    local format = {
      { "id",                       "unsigned" },
      { "statement_uuid",           "uuid",     is_nullable = true },

      { "player_id",                "number",   is_nullable = false },
      { "player_alias",             "string",   is_nullable = false }, -- This field redundancy
      { "currency_id",              "number",   is_nullable = false },


      { "channel_id",               "number",   is_nullable = false },

      { "round_id",                 "number",   is_nullable = false },
      { "round_no",                 "string",   is_nullable = false },
      { "ticket_id",                "number",   is_nullable = false },
      { "ticket_no",                "string",   is_nullable = false },

      { "total_bet_amount",         "decimal",  is_nullable = false, default = decimal.new(0.0000) }, -- Amount after exclude invalid
      { "total_bet_invalid_amount", "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "result_id",                "number",   is_nullable = true },                                 -- available only when round completed
      { "result_name",              "string",   is_nullable = true },                                 -- Name of Cock
      { "payout_amount",            "decimal",  is_nullable = true,  default = decimal.new(0.0000) }, -- Win/Lost
      { "synce_id",                 "number",   is_nullable = true,  default = 0 },                   -- The id from Billing System sync future should be UUID instead of Billing support
      { "is_sync",                  "boolean",  is_nullable = false, default = false },

      { "total_commission_amount",  "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "status_id",                "number",   is_nullable = true,  default = 1 },
      { "order",                    "number",   is_nullable = true,  default = 1 },

      { "created_by",               "number",   is_nullable = false },
      { "created_at",               "datetime", is_nullable = false },
      { "updated_by",               "number",   is_nullable = true },
      { "updated_at",               "datetime", is_nullable = true },
      { "deleted_by",               "number",   is_nullable = true },
      { "deleted_at",               "datetime", is_nullable = true },
    }

    statements_space = box.schema.create_space('statements_space', { format = format, id = 1034 })

    statements_space:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_statement_id',
      if_not_exists = true
    })

    statements_space:create_index('statement_uuid',
      {
        parts = { { 'statement_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
    statements_space:create_index('player_id',
      {
        parts = { { 'player_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
    statements_space:create_index('round_id',
      {
        parts = { { 'round_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
    statements_space:create_index('ticket_id',
      {
        parts = { { 'ticket_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

    -- box.space.statements_space:on_replace(lotto_sp.statement_trigger)
    box.sequence.seq_statement_id:set(0)
  end
end

local function init_announcement_spaces()
  local announcementSpace = box.space.announcements_space
  local seq_announcement_id = box.sequence.seq_announcement_id

  if announcementSpace == nil then
    -- Drop the existing sequence if it exists
    if seq_announcement_id ~= nil then
      seq_announcement_id:drop()
    end

    -- Create a new sequence for announcement IDs
    box.schema.sequence.create('seq_announcement_id', { start = 1 })

    -- Define the space format
    local format = {
      { "id",                       "unsigned" },                                   -- Primary key ID
      { "announcement_uuid",        "uuid",     is_nullable = false },              -- Unique identifier for the announcement
      { "announcement_desc",        "string",   is_nullable = false },              -- Description of the announcement
      { "schedule_announce",        "datetime", is_nullable = true },               -- When the announcement is scheduled
      { "schedule_announce_expire", "datetime", is_nullable = true },               -- Expiration time for the announcement
      { "announce_repeat",          "number",   is_nullable = true,  default = 0 }, -- Whether the announcement repeats (0 or 1)
      { "channel_id",               "number",   is_nullable = false },
      { "status_id",                "number",   is_nullable = false, default = 1 }, -- Status of the announcement
      { "order",                    "number",   is_nullable = true,  default = 1 }, -- Order or priority of the announcement
      { "created_at",               "datetime", is_nullable = false },              -- Record creation timestamp
      { "created_by",               "number",   is_nullable = false },              -- User ID of the creator
      { "updated_at",               "datetime", is_nullable = true },               -- Last update timestamp
      { "updated_by",               "number",   is_nullable = true },               -- User ID of the last updater
      { "deleted_at",               "datetime", is_nullable = true },               -- Deletion timestamp (if applicable)
      { "deleted_by",               "number",   is_nullable = true },               -- User ID of the deleter
    }

    -- Create the space with the specified format
    announcementSpace = box.schema.create_space('announcements_space', { format = format, id = 1035 })

    -- Create indexes
    announcementSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_announcement_id',
      if_not_exists = true,
    })
    announcementSpace:create_index('announcement_uuid', {
      parts = { { 'announcement_uuid', 'uuid' } },
      if_not_exists = true,
      unique = true,
    })

    -- Set sequence starting value
    box.sequence.seq_announcement_id:set(0)
  end

  -- Insert a default record if the space is empty
  if announcementSpace:len() == 0 then
    announcementSpace:auto_increment {
      uuid.fromstr('e2234678-c0c4-4d0b-9179-5cb5a2ece01f'), -- announcement_uuid
      'Sample announcement description',                    -- announcement_desc
      datetime.new(os.date('*t')),                          -- schedule_announce
      nil,                                                  -- schedule_announce_expire
      0,                                                    -- announce_repeat (0 for no repeat)
      1,                                                    -- channel_id
      2,                                                    -- status_id
      1,                                                    -- order
      datetime.new(os.date('*t')),                          -- created_at
      1,                                                    -- created_by
      nil,                                                  -- updated_at
      nil,                                                  -- updated_by
      nil,                                                  -- deleted_at
      nil,                                                  -- deleted_by
    }
  end
end

local function init_rounds_logs_spaces()
  local RoundsLogsSpace = box.space.rounds_logs_space
  local seq_round_log_id = box.sequence.seq_round_log_id

  if RoundsLogsSpace == nil then
    if seq_round_log_id ~= nil then
      seq_round_log_id:drop()
    end

    box.schema.sequence.create('seq_round_log_id', { start = 1 })

    local format = {
      { "id",                "unsigned" },
      { "round_log_uuid",    "uuid",     is_nullable = false },

      { "round_id",          "number",   is_nullable = true,  default = 0 },
      { "round_log_context", "string",   is_nullable = false },
      { "round_log_desc",    "string",   is_nullable = false },
      { "log_type_id",       "number",   is_nullable = false },
      { "operator",          "string",   is_nullable = false },
      { "status_id",         "number",   is_nullable = false, default = 1 },
      { "order",             "number",   is_nullable = false, default = 1 },
      { "created_by",        "number",   is_nullable = false },
      { "created_at",        "datetime", is_nullable = false },
      { "updated_by",        "number",   is_nullable = true },
      { "updated_at",        "datetime", is_nullable = true },
      { "deleted_by",        "number",   is_nullable = true },
      { "deleted_at",        "datetime", is_nullable = true },
    }

    RoundsLogsSpace = box.schema.create_space('rounds_logs_space', { format = format, id = 1036 })

    -- Create index
    RoundsLogsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_round_log_id',
      if_not_exists = true,
    })
    RoundsLogsSpace:create_index('round_log_uuid',
      {
        parts = { { 'round_log_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      })
    box.sequence.seq_round_log_id:set(0)
  end
end

local function init_players_notifications_space()
  local playerNotificationsSpace = box.space.players_notifications_space
  local seq_player_notification_id = box.sequence.seq_player_notification_id

  if playerNotificationsSpace == nil then
    if seq_player_notification_id ~= nil then
      seq_player_notification_id:drop()
    end

    box.schema.sequence.create('seq_player_notification_id', { start = 1 })

    local format = {
      { "id",                   "unsigned" },
      { "notification_uuid",    "uuid",     is_nullable = true },
      { "player_id",            "number",   is_nullable = false },

      { "context",              "string",   is_nullable = false },
      { "subject",              "string",   is_nullable = false },
      { "description",          "string",   is_nullable = true },
      { "icon_id",              "number",   is_nullable = true },
      { "notification_type_id", "number",   is_nullable = false },
      { "status_id",            "number",   is_nullable = true, default = 1 },

      { "created_by",           "number",   is_nullable = false },
      { "created_at",           "datetime", is_nullable = false },
      { "updated_by",           "number",   is_nullable = true },
      { "updated_at",           "datetime", is_nullable = true },
      { "deleted_by",           "number",   is_nullable = true },
      { "deleted_at",           "datetime", is_nullable = true },
    }

    playerNotificationsSpace = box.schema.create_space('players_notifications_space', { format = format, id = 1037 })

    playerNotificationsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_player_notification_id',
      if_not_exists = true
    })
    playerNotificationsSpace:create_index('notification_uuid',
      {
        parts = { { 'notification_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
    playerNotificationsSpace:create_index('player_id',
      {
        parts = { { 'player_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

    playerNotificationsSpace:create_index('created_at',
      {
        parts = { { 'created_at', 'datetime' } },
        if_not_exists = true,
        unique = false
      }
    )

    playerNotificationsSpace:create_index('notification_type_id',
      {
        parts = { { 'notification_type_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
  end

  if playerNotificationsSpace:len() == 0 then
    playerNotificationsSpace:auto_increment { uuid.fromstr('2f9c03cc-358f-4115-9831-736cc9817f93'), 1, "Profile", "Photo has been changed", "", nil, 1, 1, 1, datetime.parse('2023-06-13T15:05:10.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('2088fb9e-a83c-4c12-8afe-a9d21678baa9'), 2, "Deposit", "Received funds", "10 USD is received to KM001 ", 1, 2, 1, 1, datetime.parse('2023-06-13T15:06:50.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('34b2d73c-6542-4a28-8f6d-9c6a47db56ab'), 3, "Withdrawal", "Withdrawal request processed", "", nil, 3, 1, 2, datetime.parse('2023-06-14T10:00:00.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('5de834f7-59d3-41e5-b76d-1b7495366871'), 4, "Profile", "Profile updated", "", nil, 1, 1, 2, datetime.parse('2023-06-14T11:30:00.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('8ea9a0fc-2d5a-4539-bbd6-64b12c3a70df'), 5, "Bonus", "Bonus credited", "5% cashback bonus", 2, 4, 1, 3, datetime.parse('2023-06-14T12:15:00.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('a8b07d39-f75a-43cc-8660-b7090d1a4a9f'), 6, "Deposit", "Deposit successful", "50 USD added to account", 1, 2, 1, 1, datetime.parse('2023-06-14T13:00:00.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('d2351a5e-dc2a-4e56-9b0d-e984c92034f9'), 7, "Promotion", "New promotion available", "50% deposit bonus", nil, 5, 1, 4, datetime.parse('2023-06-14T13:45:00.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('d42172c8-42ad-43a0-bb6b-f0f9bc68e9f2'), 8, "Withdrawal", "Withdrawal failed", "Insufficient funds", nil, 3, 2, 2, datetime.parse('2023-06-14T14:30:00.301612345Z'), nil, nil, nil, nil }
    playerNotificationsSpace:auto_increment { uuid.fromstr('f68354a3-90a0-4998-9617-792b29e565d1'), 9, "Deposit", "Funds deposited", "100 USD received", 1, 2, 1, 5, datetime.parse('2023-06-14T15:15:00.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_users_notifications_space()
  local userNotificationsSpace = box.space.users_notifications_space
  local seq_user_notification_id = box.sequence.seq_user_notification_id

  if userNotificationsSpace == nil then
    if seq_user_notification_id ~= nil then
      seq_user_notification_id:drop()
    end

    box.schema.sequence.create('seq_user_notification_id', { start = 1 })

    local format = {
      { "id",                   "unsigned" },
      { "notification_uuid",    "uuid",     is_nullable = true },
      { "user_id",              "number",   is_nullable = false },

      { "context",              "string",   is_nullable = false },
      { "subject",              "string",   is_nullable = false },
      { "description",          "string",   is_nullable = true },
      { "icon_id",              "number",   is_nullable = true },
      { "notification_type_id", "number",   is_nullable = false },
      { "status_id",            "number",   is_nullable = true, default = 1 },

      { "created_by",           "number",   is_nullable = false },
      { "created_at",           "datetime", is_nullable = false },
      { "updated_by",           "number",   is_nullable = true },
      { "updated_at",           "datetime", is_nullable = true },
      { "deleted_by",           "number",   is_nullable = true },
      { "deleted_at",           "datetime", is_nullable = true },
    }

    userNotificationsSpace = box.schema.create_space('users_notifications_space', { format = format, id = 1038 })

    userNotificationsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_user_notification_id',
      if_not_exists = true
    })
    userNotificationsSpace:create_index('notification_uuid',
      {
        parts = { { 'notification_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )
    userNotificationsSpace:create_index('user_id',
      {
        parts = { { 'user_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

    userNotificationsSpace:create_index('created_at',
      {
        parts = { { 'created_at', 'datetime' } },
        if_not_exists = true,
        unique = false
      }
    )

    userNotificationsSpace:create_index('notification_type_id',
      {
        parts = { { 'notification_type_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
  end

  if userNotificationsSpace:len() == 0 then
    userNotificationsSpace:auto_increment { uuid.fromstr('3c3c03cc-458f-5115-9831-736cc9817f94'), 1, "Profile", "Profile photo updated", "Your profile photo has been successfully updated.", nil, 1, 1, 1, datetime.parse('2023-06-13T15:05:10.301612345Z'), nil, nil, nil, nil }
    userNotificationsSpace:auto_increment { uuid.fromstr('4d4d04dc-568f-6225-9832-746cc9827f95'), 1, "Deposit", "Funds received", "You have received 10 USD into account KM001.", 1, 2, 1, 1, datetime.parse('2023-06-12T15:06:50.301612345Z'), nil, nil, nil, nil }
    userNotificationsSpace:auto_increment { uuid.fromstr('5e5e05ed-678f-7335-9833-756cc9837f96'), 1, "Profile", "Profile picture changed", "Your profile picture has been changed successfully.", nil, 1, 1, 1, datetime.parse('2023-06-13T15:05:10.301612345Z'), nil, nil, nil, nil }
    userNotificationsSpace:auto_increment { uuid.fromstr('6f6f06fe-789f-8445-9834-766cc9847f97'), 1, "Deposit", "Deposit alert", "20 USD has been credited to account KM012.", 1, 2, 1, 1, datetime.parse('2023-06-12T15:06:50.301612345Z'), nil, nil, nil, nil }
    userNotificationsSpace:auto_increment { uuid.fromstr('4c2b7156-612e-4b92-bd3f-d5dc243f9e91'), 2, "Profile", "Profile updated", "Your profile details have been updated.", nil, 1, 1, 1, datetime.parse('2023-06-13T15:05:10.301612345Z'), nil, nil, nil, nil }
    userNotificationsSpace:auto_increment { uuid.fromstr('bf6af7fd-7bb6-4a92-bec4-d2fa37810150'), 2, "Deposit", "New funds received", "30 USD has been added to your account KM003.", 1, 2, 1, 1, datetime.parse('2023-06-12T15:06:50.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_memberships_platforms_space()
  local membershipsPlatformsSpace = box.space.memberships_platforms_space
  local seq_membership_platform_id = box.sequence.seq_membership_platform_id

  if membershipsPlatformsSpace == nil then
    if seq_membership_platform_id ~= nil then
      seq_membership_platform_id:drop()
    end

    box.schema.sequence.create('seq_membership_platform_id', { start = 1 })

    local format = {
      { "id",                       "unsigned" },
      { "membership_platform_uuid", "uuid",     is_nullable = true },
      { "platform_name",            "string",   is_nullable = false },
      { "platform_host",            "string",   is_nullable = false },
      { "platform_token",           "string",   is_nullable = false },
      { "platform_token_expire_at", "datetime", is_nullable = false },
      { "platform_extra_payload",   "string",   is_nullable = false },
      { "last_activities",          "datetime", is_nullable = false },
      { "internal_token",           "string",   is_nullable = false },
      { "internal_token_expire_at", "datetime", is_nullable = false },

      { "status_id",                "number",   is_nullable = false, default = 1 },
      { "order",                    "number",   is_nullable = false, default = 1 },

      { "created_by",               "number",   is_nullable = false },
      { "created_at",               "datetime", is_nullable = false },
      { "updated_by",               "number",   is_nullable = true },
      { "updated_at",               "datetime", is_nullable = true },
      { "deleted_by",               "number",   is_nullable = true },
      { "deleted_at",               "datetime", is_nullable = true },
    }

    membershipsPlatformsSpace = box.schema.create_space('memberships_platforms_space', { format = format, id = 1039 })

    membershipsPlatformsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_membership_platform_id',
      if_not_exists = true
    })
    membershipsPlatformsSpace:create_index('status_id',
      {
        parts = { { 'status_id', 'number' } },
        if_not_exists = true,
        unique = true
      }
    )
  end

  if membershipsPlatformsSpace:len() == 0 then
    -- Insert default records
    membershipsPlatformsSpace:auto_increment { uuid.fromstr('a2a27ec1-aa9c-41b8-8644-3f330c03cdcb'), "Mini Membership", "http://172.18.103.16:8084/api/v1/web",
      "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE3MDk5ODEwMDEsImV4cCI6MTcxMDU4NTgwMSwiaWQiOjIsInVzZXJfbmFtZSI6IlNVUEVSU0VOSU9SMDEiLCJsb2dpbl9zZXNzaW9uIjoiYjM3NTRkOWQtM2RiMC00NDBjLWJhMzUtZTYyNjE4M2Y1ZTdmIiwibWVtYmVyc2hpcF9pZCI6MSwibWVtYmVyc2hpcF9yb2xlIjoic3VwZXJzZW5pb3IiLCJjdXJyZW5jeV9pZCI6MSwibGFuZyI6ImVuIiwicm9sZV9pZCI6Mn0.utr-2pgnmVXZdeXtHogDrjYATOsBVDiYtUS26MwjgBQ",
      datetime.parse('2024-03-09T16:05:30.301612345Z'),
      "",
      datetime.parse('2023-06-02T16:05:30.301612345Z'),
      "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE3MDk5ODExODgsImV4cCI6MTcxMDU4NTk4OCwiaWQiOiI4Mzc1MWI0OC02OGYzLTQ4MDUtYTdiZC02MGFiODMxMTkzNmQiLCJ1c2VyX25hbWUiOiJpdCIsImxvZ2luX3Nlc3Npb24iOiIwZWVhMmJjYzVmZTE0MmJiYjJkYzViZmQyN2NmYzBmYyIsImxvdHRvX3V1aWQiOiJmNDk5NmMzYS02NjQ2LTQxODQtODBiZS0zZTYwYWMxMDBmZTkiLCJsb3R0b19zcGlud2hlZWxzIjpbImE5ZjlhMGM4LTVmN2ItNGMzYS04ZTNkLTlmMGM2ZDBiOGU2YSJdLCJnYW1lX3V1aWRzIjpbIjE1OTE1MmRmLTQ0OGUtNDllOC1hMGU3LTFiMjI1MzAxY2YxMiJdfQ.BiEVgPZrfB3aGVe7ZiJf2wFHkU2m4rhnaLtPbdW39AU",
      datetime.parse('2023-06-02T16:05:30.301612345Z'),
      1,
      0,
      1,
      datetime.parse('2023-06-02T16:05:30.301612345Z'), nil, nil, nil, nil }
  end
end

local function init_latests_rounds_space()
  local latestsRoundsSpace = box.space.latests_rounds_space
  local seq_latest_round_id = box.sequence.seq_latest_round_id

  if latestsRoundsSpace == nil then
    if seq_latest_round_id ~= nil then
      seq_latest_round_id:drop()
    end

    box.schema.sequence.create('seq_latest_round_id', { start = 1 })

    local format = {
      { "id",                "unsigned" },
      { "latest_round_uuid", "uuid",     is_nullable = true },
      { "latest_round_no",   "string",   is_nullable = false },
      { "channel_id",        "number",   is_nullable = false },
      { "order",             "number",   is_nullable = false, default = 1 },
      { "status_id",         "number",   is_nullable = false, default = 1 },
      { "created_by",        "number",   is_nullable = false },
      { "created_at",        "datetime", is_nullable = false },
      { "updated_by",        "number",   is_nullable = true },
      { "updated_at",        "datetime", is_nullable = true },
      { "deleted_by",        "number",   is_nullable = true },
      { "deleted_at",        "datetime", is_nullable = true },
    }

    latestsRoundsSpace = box.schema.create_space('latests_rounds_space', { format = format, id = 1040 })

    latestsRoundsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_latest_round_id',
      if_not_exists = true
    })

    latestsRoundsSpace:create_index('latest_round_uuid',
      {
        parts = { { 'latest_round_uuid', 'uuid' } },
        if_not_exists = true,
        unique = true
      }
    )

    latestsRoundsSpace:create_index('created_at',
      {
        parts = { { 'created_at', 'datetime' } },
        if_not_exists = true,
        unique = false
      }
    )
    box.sequence.seq_latest_round_id:set(0)
  end
end

local function init_rel_roles_modules_space()
  local relRolesModulesSpace = box.space.rel_roles_modules_space
  local seq_rel_role_module_id = box.sequence.seq_rel_role_module_id

  -- Check if the space already exists
  if relRolesModulesSpace == nil then
    -- Check if the sequence exists and drop it only if it's not used
    if seq_rel_role_module_id ~= nil then
      box.schema.sequence.drop('seq_rel_role_module_id')
    end

    -- Create a new sequence
    box.schema.sequence.create('seq_rel_role_module_id', { start = 1 })

    -- Define the format for the new space
    local format = {
      { "id",           "unsigned" },
      { "role_id",      "number",   is_nullable = false },
      { "module_id",    "number",   is_nullable = false },
      { "function_ids", "string",   is_nullable = false },
      { "order",        "number",   is_nullable = false, default = 1 },
      { "status_id",    "number",   is_nullable = false, default = 1 },
      { "created_by",   "number",   is_nullable = false },
      { "created_at",   "datetime", is_nullable = false },
      { "updated_by",   "number",   is_nullable = true },
      { "updated_at",   "datetime", is_nullable = true },
      { "deleted_by",   "number",   is_nullable = true },
      { "deleted_at",   "datetime", is_nullable = true },
    }

    -- Create the space
    relRolesModulesSpace = box.schema.create_space("rel_roles_modules_space", { format = format, id = 1041 })

    -- Create indexes for the space
    relRolesModulesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_rel_role_module_id',
      if_not_exists = true
    })

    relRolesModulesSpace:create_index('created_at', {
      parts = { { 'created_at', 'datetime' } },
      if_not_exists = true,
      unique = false
    })

    -- Insert default records if the space is empty
    if relRolesModulesSpace:len() == 0 then
      local created_by = 1
      local created_at = datetime.parse('2023-06-02T16:05:30.301612345Z')

      -- Role 1: All modules with function_ids "1,2,3,4"
      for module_id = 1, 14 do
        relRolesModulesSpace:auto_increment { 1, module_id, "1,2,3,4", module_id, 1, created_by, created_at, nil, nil, nil, nil }
      end

      -- Role 2: Specific function_ids per module
      local role2_data = {
        { 1, "1" }, { 2, "1,2,3" }, { 3, "1,2" }, { 4, "1,3" }, { 5, "1,3" },
        { 6, "1,2" }, { 7, "1,2,3,4" }, { 8, "1" }, { 9, "1,2,3" }, { 10, "1,2,3,4" },
        { 11, "1" }, { 12, "1,2,3" }, { 13, "1" }, { 14, "1,2,3" }
      }
      for _, entry in ipairs(role2_data) do
        local module_id, function_ids = entry[1], entry[2]
        relRolesModulesSpace:auto_increment { 2, module_id, function_ids, module_id, 1, created_by, created_at, nil, nil, nil, nil }
      end

      -- Role 3: Specific function_ids per module
      local role3_data = {
        { 1, "1" }, { 2, "1" }, { 3, "1,2" }, { 4, "1" }, { 5, "1" },
        { 6, "1,2" }, { 7, "1" }, { 8, "1" }, { 9, "1,2,3" }, { 10, "1,2,3" },
        { 11, "1" }, { 12, "1" }, { 13, "0" }, { 14, "1" }
      }
      for _, entry in ipairs(role3_data) do
        local module_id, function_ids = entry[1], entry[2]
        relRolesModulesSpace:auto_increment { 3, module_id, function_ids, module_id, 1, created_by, created_at, nil, nil, nil, nil }
      end
    end
  end
end

local function init_players_bets_limits_space() --[[ for general ratelimits on players --]]
  local players_bets_limits_space = box.space.players_bets_limits_space
  local seq_player_bet_limit_id = box.sequence.seq_player_bet_limit_id

  if players_bets_limits_space == nil then
    if seq_player_bet_limit_id ~= nil then
      seq_player_bet_limit_id:drop()
    end

    box.schema.sequence.create('seq_player_bet_limit_id', { start = 1 })

    local format = {
      { "id",                    "unsigned" },
      { "player_bet_limit_uuid", "uuid",     is_nullable = true },

      { "player_id",             "number",   is_nullable = false },
      { "max_bet_amount",        "decimal",  is_nullable = false,  default = decimal.new(0.0000) },
      { "currency_id",           "number",   is_nullable = false },

      { "status_id",             "number",   is_nullable = false, default = 1 },
      { "order",                 "number",   is_nullable = false, default = 1 },
      { "created_by",            "number",   is_nullable = false },
      { "created_at",            "datetime", is_nullable = false },
      { "updated_by",            "number",   is_nullable = true },
      { "updated_at",            "datetime", is_nullable = true },
      { "deleted_by",            "number",   is_nullable = true },
      { "deleted_at",            "datetime", is_nullable = true },
    }

    players_bets_limits_space = box.schema.create_space('players_bets_limits_space', { format = format, id = 1044 })

    players_bets_limits_space:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_player_bet_limit_id',
      if_not_exists = true
    })
    players_bets_limits_space:create_index('currency_id',
      {
        parts = { { 'currency_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )
  end
  if players_bets_limits_space:len() == 0 then
    players_bets_limits_space:auto_increment {
      uuid.fromstr('0e0667d8-cac3-4db4-97cb-1c211940c7e5'),
      1,
      decimal.new(1000000.0000),
      1, 1, 1,
      1,
      datetime.parse('2024-10-27T10:55:13Z'),
      nil, nil, nil, nil
    }
  end
end

local function init_modules_space()
  local modulesSpace = box.space.modules_space
  local seq_module_id = box.sequence.seq_module_id

  if modulesSpace == nil then
    if seq_module_id ~= nil then
      seq_module_id:drop()
    end

    box.schema.sequence.create('seq_module_id', { start = 1 })

    local format = {
      { "id",          "unsigned" },
      { "module_uuid", "uuid",     is_nullable = true },
      { "module_name", "string",   is_nullable = false },
      { "order",       "number",   is_nullable = false, default = 1 },
      { "status_id",   "number",   is_nullable = false, default = 1 },
      { "created_by",  "number",   is_nullable = false },
      { "created_at",  "datetime", is_nullable = false },
      { "updated_by",  "number",   is_nullable = true },
      { "updated_at",  "datetime", is_nullable = true },
      { "deleted_by",  "number",   is_nullable = true },
      { "deleted_at",  "datetime", is_nullable = true },
    }

    modulesSpace = box.schema.create_space('modules_space', { format = format, id = 1042 })

    modulesSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_module_id',
      if_not_exists = true
    })

    modulesSpace:create_index('module_uuid', {
      parts = { { 'module_uuid', 'uuid' } },
      if_not_exists = true,
      unique = true
    })

    modulesSpace:create_index('created_at', {
      parts = { { 'created_at', 'datetime' } },
      if_not_exists = true,
      unique = false
    })

    if modulesSpace:len() == 0 then
      insert_modules()
    end
  end
end

local function init_rel_users_channels_space()
  local relUsersChannelsSpace = box.space.rel_users_channels_space
  local seq_rel_user_channel_id = box.sequence.seq_rel_user_channel_id

  if relUsersChannelsSpace == nil then
    if seq_rel_user_channel_id ~= nil then
      box.schema.sequence.drop('seq_rel_user_channel_id')
    end

    box.schema.sequence.create('seq_rel_user_channel_id', { start = 1 })

    local format = {
      { "id",          "unsigned" },
      { "user_id",     "number",   is_nullable = false },
      { "channel_ids", "string",   is_nullable = false },
      { "order",       "number",   is_nullable = false, default = 1 },
      { "status_id",   "number",   is_nullable = false, default = 1 },
      { "created_by",  "number",   is_nullable = false },
      { "created_at",  "datetime", is_nullable = false },
      { "updated_by",  "number",   is_nullable = true },
      { "updated_at",  "datetime", is_nullable = true },
      { "deleted_by",  "number",   is_nullable = true },
      { "deleted_at",  "datetime", is_nullable = true },
    }

    -- Create the space
    relUsersChannelsSpace = box.schema.create_space("rel_users_channels_space", { format = format, id = 1043 })

    -- Create indexes for the space
    relUsersChannelsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_rel_user_channel_id',
      if_not_exists = true
    })

    relUsersChannelsSpace:create_index('created_at', {
      parts = { { 'created_at', 'datetime' } },
      if_not_exists = true,
      unique = false
    })

    -- Insert default records if the space is empty
    if relUsersChannelsSpace:len() == 0 then
      relUsersChannelsSpace:auto_increment {
        1,
        "1,2,3",
        1,
        1,
        1,
        datetime.parse('2023-06-02T16:05:30.301612345Z'), nil, nil, nil, nil }
      relUsersChannelsSpace:auto_increment {
        2,
        "1,3",
        1,
        1,
        1,
        datetime.parse('2023-06-02T16:05:30.301612345Z'), nil, nil, nil, nil }
    end
  end
end

function insert_modules()
  local modulesSpace = box.space.modules_space
  local seq_module_id = box.sequence.seq_module_id

  if not seq_module_id then
    error("Sequence 'seq_module_id' not found. Please ensure init_modules_space() has been called.")
  end

  local modules = {
    { name = "auth",              uuid = 'a2a27ec1-aa9c-41b8-8644-3f330c03cdcb' },
    { name = "announcement",      uuid = 'b3b37ec1-bb8b-42c9-9755-4f340c04dcdc' },
    { name = "bet",               uuid = 'c4c47fc2-cc9c-43da-8866-5f450d05edef' },
    { name = "channel",           uuid = 'd5d57fc3-ddad-44ea-9877-6f560e06f0f0' },
    { name = "coin",              uuid = '8196aa07-08de-412b-9161-ac5831c9321d' },
    { name = "fight_odd",         uuid = '6fbb3e7b-3690-4c00-9a1f-ab6221205c41' },
    { name = "fight_schedule",    uuid = '84b134d3-41df-49b9-b937-1b1231a213b1' },
    -- { name = "platform", uuid = '7232afaa-4ce2-4ecb-8e66-6792ed2fa6d9' },
    { name = "player",            uuid = '4a9c9860-eef3-4c61-9c69-90703f971b5a' },
    -- { name = "player_notification", uuid = 'd5927961-6ad7-493e-9833-8c6b42572cb6' },
    { name = "result",            uuid = 'e152ab19-9b32-47ff-98fd-97d1205d443a' },
    { name = "round",             uuid = 'e2933522-d870-4200-aa29-7006c2189d7b' },
    { name = "roundlog",          uuid = '1a08dc88-75cd-414f-8723-77fac625e5db' },
    { name = "user",              uuid = '46aaecec-a1ab-4098-a255-304a9df9b4fd' },
    { name = "user_audit",        uuid = '446c99ba-5deb-4244-b76d-cafba26334aa' },
    { name = "user_notification", uuid = '456b15ea-1080-49d0-ac4b-bd728546992d' }
  }

  local current_time = datetime.parse('2023-06-02T16:05:30.301612345Z')

  for _, module_info in ipairs(modules) do
    local id = seq_module_id:next()
    local module_uuid = uuid.fromstr(module_info.uuid)
    local created_by = 1
    local status_id = 1

    modulesSpace:insert({
      id,
      module_uuid,
      module_info.name,
      id,
      status_id,
      created_by,
      current_time,
      nil,
      nil,
      nil,
      nil
    })
  end
end

local function init_announcements_banners_space()
  local announcementsBannersSpace = box.space.announcements_banners_space
  local seq_announcement_banner_id = box.sequence.seq_announcement_banner_id

  if announcementsBannersSpace == nil then
    if seq_announcement_banner_id ~= nil then
      seq_announcement_banner_id:drop()
    end

    box.schema.sequence.create('seq_announcement_banner_id', { start = 1 })

    local format = {
      { "id",                       "unsigned" },
      { "announcement_uuid",        "uuid",     is_nullable = true }, 
      { "text_en",                  "string",   is_nullable = false }, 
      { "text_zh",                  "string",   is_nullable = false }, 
      { "text_km",                  "string",   is_nullable = false }, 
      { "text_color",               "string",   is_nullable = false }, 
      { "text_alignment",           "string",   is_nullable = false }, 
      { "text_padding",             "string",   is_nullable = false }, 
      { "font_size",                "string",   is_nullable = false }, 
      { "font_weight",              "string",   is_nullable = false }, 
      { "line_height",              "string",   is_nullable = false }, 
      { "padding",                  "string",   is_nullable = false }, 
      { "background_color",         "string",   is_nullable = false }, 
      { "alignment",                "string",   is_nullable = false }, 
      { "justify",                  "string",   is_nullable = false }, 
      { "background_blur",          "string",   is_nullable = false }, 
      { "background_opacity",       "string",   is_nullable = false }, 
      { "border",                   "string",   is_nullable = false }, 
      { "rounded",                  "string",   is_nullable = false }, 
      { "channel_id",               "unsigned", is_nullable = false, default = 1 },
      { "is_active",                "boolean",  is_nullable = false, default = false }, 
      { "status_id",                "number",   is_nullable = false, default = 1 },
      { "order",                    "number",   is_nullable = false, default = 1 }, 
      { "created_by",               "number",   is_nullable = false }, 
      { "created_at",               "datetime", is_nullable = false }, 
      { "updated_by",               "number",   is_nullable = true }, 
      { "updated_at",               "datetime", is_nullable = true }, 
      { "deleted_by",               "number",   is_nullable = true }, 
      { "deleted_at",               "datetime", is_nullable = true }, 
    }

    announcementsBannersSpace = box.schema.create_space('announcements_banners_space', {
      format = format,
      id = 1046
    })

    announcementsBannersSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_announcement_banner_id',
      if_not_exists = true
    })

    announcementsBannersSpace:create_index('status_id', {
      parts = { { 'status_id', 'number' } },
      if_not_exists = true,
      unique = false
    })
  end

  if announcementsBannersSpace:len() == 0 then
      -- Insert default records
      announcementsBannersSpace:auto_increment {
          uuid.fromstr('f47ac10b-58cc-4372-a567-0e02b2c3d479'), 
          "Our Website will be close in 15 minutes",                          
          "我们的网站将在15分钟后关闭",                                       
          "គេហទំព័ររបស់យើងនឹងត្រូវបិទក្នុងរយៈពេល 15 នាទី",                          
          "#FFFFFF",         -- text_color
          "text-center",        -- text_alignment
          "px-4 py-2",          -- text_padding
          "text-xl lg:text-4xl 3xl:text-6xl",            -- font_size
          "font-bold",          -- font_weight
          "leading-loose",      -- line_height
          "p-5",                -- padding
          "#E53935",         -- background_color
          "items-center",        -- alignment
          "justify-center",     -- justify
          "backdrop-blur-xl",   -- bg-blur
          "bg-opacity-80",      -- bg-opacity     
          "",                   -- border
          "rounded-md",         -- border-redius
          1,                    -- channel_id
          false,                 -- is_active
          1,                    -- status_id
          1,                    -- order
          1,                    -- created_by
          datetime.now(),       -- created_at
          nil,                  -- updated_by
          nil,                  -- updated_at
          nil,                  -- deleted_by
          nil                   -- deleted_at
      }

    announcementsBannersSpace:auto_increment {
        uuid.fromstr('0565b489-2a79-4281-9f72-bea46a3e18e0'), 
        "Our Website will be close in 30 minutes",                          
        "我们的网站将在30分钟后关闭",                                       
        "គេហទំព័ររបស់យើងនឹងត្រូវបិទក្នុងរយៈពេល 30 នាទី",                          
        "#FFFFFF",         -- text_color
        "text-center",        -- text_alignment
        "px-4 py-2",          -- text_padding
        "text-xl lg:text-4xl 3xl:text-6xl",            -- font_size
        "font-bold",          -- font_weight
        "leading-loose",      -- line_height
        "p-5",                -- padding
        "#E53935",         -- background_color
        "items-center",        -- alignment
        "justify-center",     -- justify
        "backdrop-blur-xl",   -- bg-blur
        "bg-opacity-80",      -- bg-opacity     
        "",                   -- border
        "rounded-md",         -- border-redius
        1,                    -- channel_id
        false,                 -- is_active
        1,                    -- status_id
        1,                    -- order
        1,                    -- created_by
        datetime.now(),       -- created_at
        nil,                  -- updated_by
        nil,                  -- updated_at
        nil,                  -- deleted_by
        nil                   -- deleted_at
    }
      
  end
end

local function init_bets_limits_space()
  local betsLimitsSpace = box.space.bets_limits_space
  local seq_bet_limit_id = box.sequence.seq_bet_limit_id

  if betsLimitsSpace == nil then
    if seq_bet_limit_id ~= nil then
      seq_bet_limit_id:drop()
    end

    box.schema.sequence.create('seq_bet_limit_id', { start = 1 })

    local format = {
      { "id",            "unsigned" },
      { "bet_limit_uuid", "uuid",     is_nullable = true },

      { "bet_type_id",   "number",   is_nullable = false },
      { "channel_id",    "number",   is_nullable = false, default = 1 },
      { "bet_limit",     "decimal",  is_nullable = true,  default = decimal.new(0.0000) },
      { "status_id",     "number",   is_nullable = false, default = 1 },
      { "order",         "number",   is_nullable = false, default = 1 },

      { "created_by",    "number",   is_nullable = false },
      { "created_at",    "datetime", is_nullable = false },
      { "updated_by",    "number",   is_nullable = true },
      { "updated_at",    "datetime", is_nullable = true },
      { "deleted_by",    "number",   is_nullable = true },
      { "deleted_at",    "datetime", is_nullable = true },
    }

    betsLimitsSpace = box.schema.create_space('bets_limits_space', { format = format, id = 1047 })

    betsLimitsSpace:create_index('id', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_bet_limit_id',
      if_not_exists = true
    })

    betsLimitsSpace:create_index('bet_limit_uuid', {
      parts = { { 'bet_limit_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })

    betsLimitsSpace:create_index('channel_id',
      {
        parts = { { 'channel_id', 'number' } },
        if_not_exists = true,
        unique = false
      }
    )

  end

  if betsLimitsSpace:len() == 0 then
    betsLimitsSpace:auto_increment { uuid.fromstr('a1b2c3d4-e5f6-4789-b0c1-d2e3f4a5b6c7'), 1, 1, decimal.new(5000000), 1, 1, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('b2c3d4e5-f6a7-4890-b1c2-d3e4f5a6b7c8'), 2, 1, decimal.new(5000000), 1, 2, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('c3d4e5f6-a7b8-4901-b2c3-d4e5f6a7b8c9'), 3, 1, decimal.new(5000000), 1, 3, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('d4e5f6a7-b8c9-4012-b3c4-d5e6f7a8b9c0'), 4, 1, decimal.new(5000000), 1, 4, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }

    -- Insert default records for channel 2
    betsLimitsSpace:auto_increment { uuid.fromstr('e5f6a7b8-c9d0-4123-b4c5-d6e7f8a9b0c1'), 5, 2, decimal.new(5000000), 1, 5, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('f6a7b8c9-d0e1-4234-b5c6-d7e8f9a0b1c2'), 6, 2, decimal.new(5000000), 1, 6, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('a7b8c9d0-e1f2-4345-b6c7-d8e9f0a1b2c3'), 7, 2, decimal.new(5000000), 1, 7, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('b8c9d0e1-f2a3-4456-b7c8-d9e0f1a2b3c4'), 8, 2, decimal.new(5000000), 1, 8, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }

    -- Insert default records for channel 3
    betsLimitsSpace:auto_increment { uuid.fromstr('c9d0e1f2-a3b4-4567-b8c9-e0f1a2b3c4d5'), 9, 3, decimal.new(5000000), 1, 9, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('d0e1f2a3-b4c5-4678-b9c0-e1f2a3b4c5d6'), 10, 3, decimal.new(5000000), 1, 10, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('e1f2a3b4-c5d6-4789-b0c1-f2a3b4c5d6e7'), 11, 3, decimal.new(5000000), 1, 11, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }
    betsLimitsSpace:auto_increment { uuid.fromstr('f2a3b4c5-d6e7-4890-b1c2-a3b4c5d6e7f8'), 12, 3, decimal.new(5000000), 1, 12, 1, datetime.parse('2025-02-18T09:55:30.301612345Z'), nil, nil, nil, nil }

  end
end

local function init_cards_space()

  local cardsSpace = box.space.cards_space
  local seq_card_id = box.sequence.seq_card_id

  if cardsSpace == nil then

    if seq_card_id ~= nil then
      seq_card_id:drop()
    end

    box.schema.sequence.create('seq_card_id', { start = 1 })

    local format = {
      { "id",           "unsigned" },
      { "card_uuid",    "uuid",     is_nullable = true },

      { "card_suit_id", "number",   is_nullable = false }, -- 1♥ 2♦ 3♣ 4♠
      { "card_name",    "string",   is_nullable = false },
      { "card_number",  "number",   is_nullable = false },
      { "card_value",   "number",   is_nullable = false }, -- Pokdeng value
      { "card_type_id", "number",   is_nullable = false }, -- 1=normal 2=face

      { "card_image",   "string",   is_nullable = true },

      { "status_id",    "number",   is_nullable = false, default = 1 },
      { "order",        "number",   is_nullable = false, default = 1 },

      { "created_by",   "number",   is_nullable = true },
      { "created_at",   "datetime", is_nullable = false },
      { "updated_by",   "number",   is_nullable = true },
      { "updated_at",   "datetime", is_nullable = true },
      { "deleted_by",   "number",   is_nullable = true },
      { "deleted_at",   "datetime", is_nullable = true },
    }

    cardsSpace = box.schema.create_space('cards_space', {
      format = format,
      id = 1050,
      if_not_exists = true
    })

    cardsSpace:create_index('primary', {
      parts = { { 'id', 'unsigned' } },
      sequence = 'seq_card_id',
      if_not_exists = true
    })

    cardsSpace:create_index('card_uuid', {
      parts = { { 'card_uuid', 'uuid' } },
      unique = true,
      if_not_exists = true
    })

    cardsSpace:create_index('card_number', {
      parts = { { 'card_number', 'number' } },
      unique = false,
      if_not_exists = true
    })

    cardsSpace:create_index('card_suit_id', {
      parts = { { 'card_suit_id', 'number' } },
      unique = false,
      if_not_exists = true
    })
  end

  -- ✅ Seed only once
  if cardsSpace:len() > 0 then
    return
  end

  local now = datetime.now()

  local function insert(card)
    cardsSpace:auto_increment(card)
  end

  -- ✅ Suits
  local suits = {
    { id = 1, name = 'heart' },
    { id = 2, name = 'diamond' },
    { id = 3, name = 'club' },
    { id = 4, name = 'spade' },
  }

  -- ✅ Numbers & Pokdeng values
  local cards = {
    [1]  = { 'ace',   1 },
    [2]  = { 'two',   2 },
    [3]  = { 'three', 3 },
    [4]  = { 'four',  4 },
    [5]  = { 'five',  5 },
    [6]  = { 'six',   6 },
    [7]  = { 'seven', 7 },
    [8]  = { 'eight', 8 },
    [9]  = { 'nine',  9 },
    [10] = { 'ten',   10 },
    [11] = { 'jack',  10 },
    [12] = { 'queen', 10 },
    [13] = { 'king',  10 },
  }

  for number, info in pairs(cards) do
    for _, suit in ipairs(suits) do
      local card_type_id = number >= 11 and 2 or 1

      insert {
        uuid.new(),
        suit.id,
        info[1],
        number,
        info[2],
        card_type_id,
        string.format('%d-%s.webp', number, suit.name),
        1,
        1,
        1,
        now,
        nil, nil, nil, nil
      }
    end
  end
end

local function init_rooms_space()
  local space_name = 'rooms_space'
  local seq_name = 'seq_room_id'

  if box.space[space_name] ~= nil then
    return
  end

  if box.sequence[seq_name] ~= nil then
    box.sequence[seq_name]:drop()
  end

  box.schema.sequence.create(seq_name, { start = 1 })

  local format = {
    { "id",           "unsigned" },
    { "room_uuid",    "uuid",     is_nullable = true },
    { "room_code",    "string",   is_nullable = false },
    { "room_name",    "string",   is_nullable = false },
    { "currency_id",  "number",   is_nullable = false },
    { "min_bet",      "decimal",  is_nullable = false },
    { "max_bet",      "decimal",  is_nullable = false },
    { "status_id",    "number",   is_nullable = false, default = 1 },
    { "created_by",   "number",   is_nullable = false },
    { "created_at",   "datetime", is_nullable = false },
    { "updated_by",   "number",   is_nullable = true },
    { "updated_at",   "datetime", is_nullable = true },
    { "deleted_by",   "number",   is_nullable = true },
    { "deleted_at",   "datetime", is_nullable = true },
  }

  local s = box.schema.create_space(space_name, {
    id = 1051,
    format = format,
  })

  s:create_index('pk', { parts = { { 'id', 'unsigned' } }, sequence = seq_name })
  s:create_index('room_uuid', { parts = { { 'room_uuid', 'uuid' } }, unique = true })
  s:create_index('room_code', { parts = { { 'room_code', 'string' } }, unique = true })
  box.sequence[seq_name]:set(0)

  -- INSERT DEFAULT ROOMS
  local provinces = {
    "Phnom Penh", "Kandal", "Siem Reap", "Battambang", "Banteay Meanchey",
    "Kampong Cham", "Kampong Chhnang", "Kampong Speu", "Kampong Thom",
    "Kampot", "Kep", "Koh Kong", "Kratie", "Mondulkiri", "Oddar Meanchey",
    "Pailin", "Preah Vihear", "Prey Veng", "Pursat", "Ratanakiri",
    "Stung Treng", "Svay Rieng", "Takeo", "Tbong Khmum"
  }

  for _, name in ipairs(provinces) do
    local code = string.upper(string.gsub(name, "%s+", "")) -- uppercase + trim
    s:insert({
      nil,
      uuid.new(),
      code,
      name,
      1,                   -- currency_id
      decimal.new(1000),   -- min_bet
      decimal.new(1000000),-- max_bet
      1,
      0,
      datetime.now(),
      nil,
      nil,
      nil,
      nil,
    })
  end
end

local function init_tables_space()
  local space_name = 'tables_space'
  local seq_name = 'seq_table_id'

  if box.space[space_name] == nil then
    if box.sequence[seq_name] ~= nil then
      box.sequence[seq_name]:drop()
    end
    box.schema.sequence.create(seq_name, { start = 1 })

    local format = {
      { "id",            "unsigned" },
      { "table_uuid",    "uuid", is_nullable = true },
      { "room_id",       "number", is_nullable = false },
      { "table_code",    "string", is_nullable = false },
      { "table_name",    "string", is_nullable = false },
      { "max_players",   "number", is_nullable = false, default = 6 },
      { "currency_id",   "number", is_nullable = false },
      { "min_bet",       "decimal", is_nullable = false },
      { "max_bet",       "decimal", is_nullable = false },
      { "status_id",     "number", is_nullable = false, default = 1 },
      { "created_by",    "number", is_nullable = false },
      { "created_at",    "datetime", is_nullable = false },
      { "updated_by",    "number", is_nullable = true },
      { "updated_at",    "datetime", is_nullable = true },
      { "deleted_by",    "number", is_nullable = true },
      { "deleted_at",    "datetime", is_nullable = true },
    }

    local s = box.schema.create_space(space_name, {
      id = 1052,
      format = format,
    })

    s:create_index('pk', { parts = { { 'id', 'unsigned' } }, sequence = seq_name })
    s:create_index('table_uuid', { parts = { { 'table_uuid', 'uuid' } }, unique = true })
    s:create_index('room_id', { parts = { { 'room_id', 'number' } }, unique = false })
    s:create_index('room_table_code', {
      parts = { { 'room_id', 'number' }, { 'table_code', 'string' } },
      unique = true,
    })

    box.sequence[seq_name]:set(0)
  end

  -- CREATE 5 TABLES PER ROOM
  local rooms = box.space.rooms_space:select{} -- select all rooms
  for _, room in ipairs(rooms) do
    local room_id = room.id
    local room_code = room.room_code
    local room_name = room.room_name

    for t = 1, 5 do
      local suffix = string.format("%03d", t) -- 001, 002...
      local table_code = room_code .. suffix
      local table_name = room_name .. " " .. suffix

      box.space.tables_space:insert({
        nil,
        uuid.new(),
        room_id,
        table_code,
        table_name,
        6,                   -- max players
        1,                   -- currency_id
        decimal.new(1000),   -- min_bet
        decimal.new(1000000),-- max_bet
        1,
        0,
        datetime.now(),
        nil,
        nil,
        nil,
        nil,
      })
    end
  end
end

init_user_spaces()
init_users_roles_space()
init_users_audits_spaces()
init_admin_admin_menus_space()
init_users_menus_spaces()
init_currencies_space()
init_languages_space()
init_exchange_rates_space()
init_currencies_defaults_rates_space()
init_channels_space()
init_rounds_space()
init_results_space()
init_cocks_space()
init_fights_schedules_space()
init_fights_schedules_details_space()
init_fights_odds_space()
init_bets_types_space()
init_defaults_ratelimits_space()
init_coins_assets_space()
init_coins_space()
init_member_spaces()
init_members_balances_space()
init_players_roles_space()
init_notifications_types_space()
init_countries_space()
init_players_avatars_space()
init_settings_music_space()
init_settings_volume_space()
init_settings_voice_space()
init_players_ratelimits_space()
init_players_audits_spaces()
init_tickets_space()
init_bets_space()
init_statements_space()
init_announcement_spaces()
init_rounds_logs_spaces()
init_players_notifications_space()
init_users_notifications_space()
init_memberships_platforms_space()
init_latests_rounds_space()
init_rel_roles_modules_space()
init_modules_space()
init_rel_users_channels_space()
init_announcements_banners_space()
init_bets_limits_space()
init_players_bets_limits_space()
init_cards_space()
init_rooms_space()
init_tables_space()



