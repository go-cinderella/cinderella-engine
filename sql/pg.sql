create table if not exists act_hi_identitylink
(
    id_                  varchar(64) not null
        primary key,
    group_id_            varchar(255),
    type_                varchar(255),
    user_id_             varchar(255),
    task_id_             varchar(64),
    create_time_         timestamp,
    proc_inst_id_        varchar(64),
    scope_id_            varchar(255),
    sub_scope_id_        varchar(255),
    scope_type_          varchar(255),
    scope_definition_id_ varchar(255)
);

create index if not exists act_idx_hi_ident_lnk_user
    on act_hi_identitylink (user_id_);

create index if not exists act_idx_hi_ident_lnk_scope
    on act_hi_identitylink (scope_id_, scope_type_);

create index if not exists act_idx_hi_ident_lnk_sub_scope
    on act_hi_identitylink (sub_scope_id_, scope_type_);

create index if not exists act_idx_hi_ident_lnk_scope_def
    on act_hi_identitylink (scope_definition_id_, scope_type_);

create index if not exists act_idx_hi_ident_lnk_task
    on act_hi_identitylink (task_id_);

create index if not exists act_idx_hi_ident_lnk_procinst
    on act_hi_identitylink (proc_inst_id_);

create table if not exists act_hi_taskinst
(
    id_                       varchar(64) not null
        primary key,
    rev_                      integer      default 1,
    proc_def_id_              varchar(64),
    task_def_id_              varchar(64),
    task_def_key_             varchar(255),
    proc_inst_id_             varchar(64),
    execution_id_             varchar(64),
    scope_id_                 varchar(255),
    sub_scope_id_             varchar(255),
    scope_type_               varchar(255),
    scope_definition_id_      varchar(255),
    propagated_stage_inst_id_ varchar(255),
    name_                     varchar(255),
    parent_task_id_           varchar(64),
    description_              varchar(4000),
    owner_                    varchar(255),
    assignee_                 varchar(255),
    start_time_               timestamp   not null,
    claim_time_               timestamp,
    end_time_                 timestamp,
    duration_                 bigint,
    delete_reason_            varchar(4000),
    priority_                 integer,
    due_date_                 timestamp,
    form_key_                 varchar(255),
    category_                 varchar(255),
    tenant_id_                varchar(255) default ''::varchar,
    last_updated_time_        timestamp
);

create index if not exists act_idx_hi_task_scope
    on act_hi_taskinst (scope_id_, scope_type_);

create index if not exists act_idx_hi_task_sub_scope
    on act_hi_taskinst (sub_scope_id_, scope_type_);

create index if not exists act_idx_hi_task_scope_def
    on act_hi_taskinst (scope_definition_id_, scope_type_);

create index if not exists act_idx_hi_task_inst_procinst
    on act_hi_taskinst (proc_inst_id_);

create table if not exists act_hi_varinst
(
    id_                varchar(64)  not null
        primary key,
    rev_               integer default 1,
    proc_inst_id_      varchar(64),
    execution_id_      varchar(64),
    task_id_           varchar(64),
    name_              varchar(255) not null,
    var_type_          varchar(100),
    scope_id_          varchar(255),
    sub_scope_id_      varchar(255),
    scope_type_        varchar(255),
    bytearray_id_      varchar(64),
    double_            double precision,
    long_              bigint,
    text_              varchar(4000),
    text2_             varchar(4000),
    create_time_       timestamp,
    last_updated_time_ timestamp
);

create index if not exists act_idx_hi_procvar_name_type
    on act_hi_varinst (name_, var_type_);

create index if not exists act_idx_hi_var_scope_id_type
    on act_hi_varinst (scope_id_, scope_type_);

create index if not exists act_idx_hi_var_sub_id_type
    on act_hi_varinst (sub_scope_id_, scope_type_);

create index if not exists act_idx_hi_procvar_proc_inst
    on act_hi_varinst (proc_inst_id_);

create index if not exists act_idx_hi_procvar_task_id
    on act_hi_varinst (task_id_);

create index if not exists act_idx_hi_procvar_exe
    on act_hi_varinst (execution_id_);

create unique index if not exists act_hi_varinst__uniq
    on act_hi_varinst (name_, proc_inst_id_);

create table if not exists act_re_deployment
(
    id_                   varchar(64) not null
        primary key,
    name_                 varchar(255),
    category_             varchar(255),
    key_                  varchar(255),
    tenant_id_            varchar(255) default ''::varchar,
    deploy_time_          timestamp,
    derived_from_         varchar(64),
    derived_from_root_    varchar(64),
    parent_deployment_id_ varchar(255),
    engine_version_       varchar(255),
    process_id_           text
);

create table if not exists act_ge_bytearray
(
    id_            varchar(64) not null
        primary key,
    rev_           integer,
    name_          varchar(255),
    deployment_id_ varchar(64)
        constraint act_fk_bytearr_depl
            references act_re_deployment,
    bytes_         bytea,
    generated_     boolean
);

create index if not exists act_idx_bytear_depl
    on act_ge_bytearray (deployment_id_);

create table if not exists act_re_procdef
(
    id_                     varchar(64)            not null
        primary key,
    rev_                    integer,
    category_               varchar(255),
    name_                   varchar(255),
    key_                    varchar(255)           not null,
    version_                integer                not null,
    deployment_id_          varchar(64),
    resource_name_          varchar(4000),
    dgrm_resource_name_     varchar(4000),
    description_            varchar(4000),
    has_start_form_key_     boolean,
    has_graphical_notation_ boolean,
    suspension_state_       integer,
    tenant_id_              varchar(255) default ''::varchar,
    derived_from_           varchar(64),
    derived_from_root_      varchar(64),
    derived_version_        integer      default 0 not null,
    engine_version_         varchar(255),
    process_id_             varchar(255),
    constraint act_uniq_procdef
        unique (key_, version_, derived_version_, tenant_id_)
);

create table if not exists act_ru_execution
(
    id_                        varchar(64) not null
        primary key,
    rev_                       integer,
    proc_inst_id_              varchar(64)
        constraint act_fk_exe_procinst
            references act_ru_execution,
    business_key_              varchar(255),
    parent_id_                 varchar(64)
        constraint act_fk_exe_parent
            references act_ru_execution,
    proc_def_id_               varchar(64)
        constraint act_fk_exe_procdef
            references act_re_procdef,
    super_exec_                varchar(64)
        constraint act_fk_exe_super
            references act_ru_execution,
    root_proc_inst_id_         varchar(64),
    act_id_                    varchar(255),
    is_active_                 boolean,
    is_concurrent_             boolean,
    is_scope_                  boolean,
    is_event_scope_            boolean,
    is_mi_root_                boolean,
    suspension_state_          integer,
    cached_ent_state_          integer,
    tenant_id_                 varchar(255) default ''::varchar,
    name_                      varchar(255),
    start_act_id_              varchar(255),
    start_time_                timestamp,
    start_user_id_             varchar(255),
    lock_time_                 timestamp,
    lock_owner_                varchar(255),
    is_count_enabled_          boolean,
    evt_subscr_count_          integer,
    task_count_                integer,
    job_count_                 integer,
    timer_job_count_           integer,
    susp_job_count_            integer,
    deadletter_job_count_      integer,
    external_worker_job_count_ integer,
    var_count_                 integer,
    id_link_count_             integer,
    callback_id_               varchar(255),
    callback_type_             varchar(255),
    reference_id_              varchar(255),
    reference_type_            varchar(255),
    propagated_stage_inst_id_  varchar(255),
    business_status_           varchar(255)
);

create table if not exists act_ru_task
(
    id_                       varchar(64) not null
        primary key,
    rev_                      integer,
    execution_id_             varchar(64)
        constraint act_fk_task_exe
            references act_ru_execution,
    proc_inst_id_             varchar(64)
        constraint act_fk_task_procinst
            references act_ru_execution,
    proc_def_id_              varchar(64)
        constraint act_fk_task_procdef
            references act_re_procdef,
    task_def_id_              varchar(64),
    scope_id_                 varchar(255),
    sub_scope_id_             varchar(255),
    scope_type_               varchar(255),
    scope_definition_id_      varchar(255),
    propagated_stage_inst_id_ varchar(255),
    name_                     varchar(255),
    parent_task_id_           varchar(64),
    description_              varchar(4000),
    task_def_key_             varchar(255),
    owner_                    varchar(255),
    assignee_                 varchar(255),
    delegation_               varchar(64),
    priority_                 integer,
    create_time_              timestamp,
    due_date_                 timestamp,
    category_                 varchar(255),
    suspension_state_         integer,
    tenant_id_                varchar(255) default ''::varchar,
    form_key_                 varchar(255),
    claim_time_               timestamp,
    is_count_enabled_         boolean,
    var_count_                integer,
    id_link_count_            integer,
    sub_task_count_           integer
);

create table if not exists act_ru_identitylink
(
    id_                  varchar(64) not null
        primary key,
    rev_                 integer,
    group_id_            varchar(255),
    type_                varchar(255),
    user_id_             varchar(255),
    task_id_             varchar(64)
        constraint act_fk_tskass_task
            references act_ru_task,
    proc_inst_id_        varchar(64)
        constraint act_fk_idl_procinst
            references act_ru_execution,
    proc_def_id_         varchar(64)
        constraint act_fk_athrz_procedef
            references act_re_procdef,
    scope_id_            varchar(255),
    sub_scope_id_        varchar(255),
    scope_type_          varchar(255),
    scope_definition_id_ varchar(255)
);

create index if not exists act_idx_ident_lnk_user
    on act_ru_identitylink (user_id_);

create index if not exists act_idx_ident_lnk_group
    on act_ru_identitylink (group_id_);

create index if not exists act_idx_ident_lnk_scope
    on act_ru_identitylink (scope_id_, scope_type_);

create index if not exists act_idx_ident_lnk_sub_scope
    on act_ru_identitylink (sub_scope_id_, scope_type_);

create index if not exists act_idx_ident_lnk_scope_def
    on act_ru_identitylink (scope_definition_id_, scope_type_);

create index if not exists act_idx_tskass_task
    on act_ru_identitylink (task_id_);

create index if not exists act_idx_athrz_procedef
    on act_ru_identitylink (proc_def_id_);

create index if not exists act_idx_idl_procinst
    on act_ru_identitylink (proc_inst_id_);

create index if not exists act_idx_task_create
    on act_ru_task (create_time_);

create index if not exists act_idx_task_scope
    on act_ru_task (scope_id_, scope_type_);

create index if not exists act_idx_task_sub_scope
    on act_ru_task (sub_scope_id_, scope_type_);

create index if not exists act_idx_task_scope_def
    on act_ru_task (scope_definition_id_, scope_type_);

create index if not exists act_idx_task_exec
    on act_ru_task (execution_id_);

create index if not exists act_idx_task_procinst
    on act_ru_task (proc_inst_id_);

create index if not exists act_idx_task_procdef
    on act_ru_task (proc_def_id_);

create table if not exists act_ru_variable
(
    id_           varchar(64)  not null
        primary key,
    rev_          integer,
    type_         varchar(255) not null,
    name_         varchar(255) not null,
    execution_id_ varchar(64)
        constraint act_fk_var_exe
            references act_ru_execution,
    proc_inst_id_ varchar(64)
        constraint act_fk_var_procinst
            references act_ru_execution,
    task_id_      varchar(64),
    scope_id_     varchar(255),
    sub_scope_id_ varchar(255),
    scope_type_   varchar(255),
    bytearray_id_ varchar(64)
        constraint act_fk_var_bytearray
            references act_ge_bytearray,
    double_       double precision,
    long_         bigint,
    text_         varchar(4000),
    text2_        varchar(4000)
);

create index if not exists act_idx_ru_var_scope_id_type
    on act_ru_variable (scope_id_, scope_type_);

create index if not exists act_idx_ru_var_sub_id_type
    on act_ru_variable (sub_scope_id_, scope_type_);

create index if not exists act_idx_var_bytearray
    on act_ru_variable (bytearray_id_);

create index if not exists act_idx_variable_task_id
    on act_ru_variable (task_id_);

create index if not exists act_idx_var_exe
    on act_ru_variable (execution_id_);

create index if not exists act_idx_var_procinst
    on act_ru_variable (proc_inst_id_);

create unique index if not exists act_ru_variable__uniq
    on act_ru_variable (name_, proc_inst_id_);

create index if not exists act_idx_exec_buskey
    on act_ru_execution (business_key_);

create index if not exists act_idx_exe_root
    on act_ru_execution (root_proc_inst_id_);

create index if not exists act_idx_exec_ref_id_
    on act_ru_execution (reference_id_);

create index if not exists act_idx_exe_procinst
    on act_ru_execution (proc_inst_id_);

create index if not exists act_idx_exe_parent
    on act_ru_execution (parent_id_);

create index if not exists act_idx_exe_super
    on act_ru_execution (super_exec_);

create index if not exists act_idx_exe_procdef
    on act_ru_execution (proc_def_id_);

create table if not exists act_hi_procinst
(
    id_                        varchar(64) not null
        primary key,
    rev_                       integer      default 1,
    proc_inst_id_              varchar(64) not null
        unique,
    business_key_              varchar(255),
    proc_def_id_               varchar(64) not null,
    start_time_                timestamp   not null,
    end_time_                  timestamp,
    duration_                  bigint,
    start_user_id_             varchar(255),
    start_act_id_              varchar(255),
    end_act_id_                varchar(255),
    super_process_instance_id_ varchar(64),
    delete_reason_             varchar(4000),
    tenant_id_                 varchar(255) default ''::varchar,
    name_                      varchar(255),
    callback_id_               varchar(255),
    callback_type_             varchar(255),
    reference_id_              varchar(255),
    reference_type_            varchar(255),
    propagated_stage_inst_id_  varchar(255),
    business_status_           varchar(255)
);

create index if not exists act_idx_hi_pro_inst_end
    on act_hi_procinst (end_time_);

create index if not exists act_idx_hi_pro_i_buskey
    on act_hi_procinst (business_key_);

create index if not exists act_idx_hi_pro_super_procinst
    on act_hi_procinst (super_process_instance_id_);

create table if not exists act_hi_actinst
(
    id_                 varchar(64)  not null
        primary key,
    rev_                integer      default 1,
    proc_def_id_        varchar(64)  not null,
    proc_inst_id_       varchar(64)  not null,
    execution_id_       varchar(64)  not null,
    act_id_             varchar(255) not null,
    task_id_            varchar(64),
    call_proc_inst_id_  varchar(64),
    act_name_           varchar(255),
    act_type_           varchar(255) not null,
    assignee_           varchar(255),
    start_time_         timestamp    not null,
    end_time_           timestamp,
    transaction_order_  integer,
    duration_           bigint,
    delete_reason_      varchar(4000),
    tenant_id_          varchar(255) default ''::varchar,
    business_result_    text,
    business_parameter_ text
);

create index if not exists act_idx_hi_act_inst_start
    on act_hi_actinst (start_time_);

create index if not exists act_idx_hi_act_inst_end
    on act_hi_actinst (end_time_);

create index if not exists act_idx_hi_act_inst_procinst
    on act_hi_actinst (proc_inst_id_, act_id_);

create index if not exists act_idx_hi_act_inst_exec
    on act_hi_actinst (execution_id_, act_id_);

