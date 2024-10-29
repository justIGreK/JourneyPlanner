// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/signIn": {
            "post": {
                "description": "Authorization to the account",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "SignIn",
                "parameters": [
                    {
                        "type": "string",
                        "description": "your login or email",
                        "name": "option",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "your password",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/auth/singUp": {
            "post": {
                "description": "Create account",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "SignUp",
                "parameters": [
                    {
                        "type": "string",
                        "description": "your login",
                        "name": "login",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "your password",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "your email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/add": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create new group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "AddGroup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of group",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "multi",
                        "description": "by adding logins you will automatically invite this users",
                        "name": "invites",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/groups/ban": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Kick and ban member from group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blacklist"
                ],
                "summary": "BanMember",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "login of member",
                        "name": "memberLogin",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/blacklist": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get blacklist of group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blacklist"
                ],
                "summary": "Get group blacklist",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/declineinvite": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Decline invite",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invites"
                ],
                "summary": "Decline invite",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of invite",
                        "name": "invite_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/delete": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete group by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "DeleteGroup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/getgroupinfo": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get full info about group you are a member of",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "GetGroupInfo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/getlist": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get a list of all the groups you are a member of",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "GetGroups",
                "responses": {}
            }
        },
        "/groups/givelead": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Give another member of group leader role",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "GiveLeaderRole",
                "parameters": [
                    {
                        "type": "string",
                        "description": "member login",
                        "name": "user_login",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/invite": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Invite user to group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invites"
                ],
                "summary": "Invite user to group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "invited user",
                        "name": "user_login",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/invitelist": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get your list of invites",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invites"
                ],
                "summary": "Get invite list",
                "responses": {}
            }
        },
        "/groups/leaveGroup": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Leave from group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "LeaveFromGroup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/groups/unban": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Unban member in group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blacklist"
                ],
                "summary": "UnbanMember",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "login of member",
                        "name": "memberLogin",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/polls/add": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create new poll",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "polls"
                ],
                "summary": "CreatePoll",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "groupID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "title of poll",
                        "name": "title",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "first option ",
                        "name": "firstOption",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "second option",
                        "name": "sercondOption",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "duration of poll in minutes",
                        "name": "duration",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/polls/close": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Close poll for voting",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "polls"
                ],
                "summary": "Close Poll",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "groupID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "id of poll",
                        "name": "pollID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/polls/delete": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete poll by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "polls"
                ],
                "summary": "Delete Poll",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "groupID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "id of poll",
                        "name": "pollID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/polls/getlist": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get list of polls",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "polls"
                ],
                "summary": "GetPolls",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "groupID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/polls/vote": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Vote for poll option",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "polls"
                ],
                "summary": "Vote Poll",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of group",
                        "name": "groupID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "id of poll",
                        "name": "pollID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "firstOption",
                            "secondOption"
                        ],
                        "type": "string",
                        "description": "vote option",
                        "name": "option",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/tasks/add": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create new task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "AddTask",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Task Details",
                        "name": "title",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "2024-10-21",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "14:00",
                        "name": "start_time",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 0,
                        "name": "days",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 2,
                        "name": "hours",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 30,
                        "name": "minutes",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/tasks/delete": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete existing task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "DeleteTask",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Id of group",
                        "name": "task_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/tasks/getlist": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create new task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "GetTasks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/tasks/update": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "update existing task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "UpdateTask",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of group",
                        "name": "group_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "task id",
                        "name": "task_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Task Details",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "2024-10-21",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "14:00",
                        "name": "start_time",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 0,
                        "name": "days",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 2,
                        "name": "hours",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 30,
                        "name": "minutes",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Journer Planner",
	Description:      "Application for planning your journey",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
