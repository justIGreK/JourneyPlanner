# Journey Planner
## Overview

Journey Planner is a group travel planning tool built in Go, designed to simplify collaborative trip planning. It allows users to organize group activities, manage tasks, and facilitate voting within the group. The system supports seamless communication through WebSocket-based group chats, ensuring only authorized group members can connect and participate.

## Features
- **Group Management**: Users can create groups, add friends, assign roles (e.g., leaders, moderators), and manage memberships.
- **Invitations**: Group members can generate invite links to add new members, with security measures to ensure only authorized users can create invites.
- **Task Management**: Group leaders can create, modify, and delete tasks with structured timelines, ensuring no overlapping tasks.
- **Voting**: A built-in voting system lets group members vote on tasks or ideas, with options to retrieve open and closed polls.
- **WebSocket Chat**: Real-time group chats allow members to discuss and coordinate plans, with enforced group membership and chat history retrieval upon connection.
- **Role-Based Permissions**: Group leaders can edit the group composition, kicking out of the group, simultaneously blacklisting the user, or unbanning.
- **MongoDB Integration**: Data is stored and managed with MongoDB, providing robust handling of groups, tasks, and voting data.

## Project Structure

- **cmd/handler**: Contains all application endpoints, handling HTTP requests for different functionalities.
- **cmd/handler/ws** Manages WebSocket connections, including message handling and disconnection logic.
- **internal/service**: Implements services and business logic, centralizing application functionality.
- **internal/repository** Manages database operations, acting as an interface between the database and services.
- **internal/models** Defines core data structures used throughout the application.

## Installation and Running

1. Clone the repository:

   ```bash
    git clone https://github.com/justIGreK/JourneyPlanner
2. Build and run the service using Docker Compose:
   ```bash
    docker-compose up --build
## Usage 
### Users
- `POST` /auth/signIn - SignIn: Logs in a user to the application.
- `POST` /auth/signUp - SignUp: Registers a new user.
### Groups
- `POST` /groups/add - AddGroup: Creates a new group.
- `DELETE` /groups/delete - DeleteGroup: Removes an existing group.
- `GET` /groups/getgroupinfo - GetGroupInfo: Retrieves information about a specific group.
- `GET` /groups/getlist - GetGroups: Retrieves a list of groups that the user belongs to.
- `PUT` /groups/givelead - GiveLeaderRole: Assigns the leader role to a specified member.
- `POST` /groups/leaveGroup - LeaveFromGroup: Allows a user to leave a group.
### Blacklist Management
- `PUT` /groups/ban - BanMember: Bans a member from the group.
- `GET` /groups/blacklist - Get group blacklist: Retrieves the blacklist of banned members.
- `PUT` /groups/unban - UnbanMember: Removes a member from the blacklist.
### Invitations
- `POST` /groups/declineinvite - Decline Invite: Declines an invitation to join a group.
- `POST` /groups/invite - Invite user to group: Sends an invitation to a user to join a group.
- `GET` /groups/invitelist - Get invite list: Retrieves the list of pending group invitations.
### Polls
- `POST` /polls/add - CreatePoll: Creates a new poll within a group.
- `PUT` /polls/close - Close Poll: Closes an active poll.
- `DELETE` /polls/delete - Delete Poll: Deletes an existing poll.
- `GET` /polls/getlist - GetPolls: Retrieves a list of open and closed polls.
- `PUT` /polls/vote - Vote Poll: Casts a vote in a poll.
### Tasks
- `POST` /tasks/add - AddTask: Adds a new task to a group.
- `DELETE` /tasks/delete - DeleteTask: Removes an existing task.
- `GET` /tasks/getlist - GetTasks: Retrieves a list of tasks in a group.
- `PUT` /tasks/update - UpdateTask: Updates the details of an existing task.
#### Testing Functionality
For most endpoints, an authorization token is required. This token is provided upon a successful login and must be included in the **Authorization** header with the **Bearer** prefix.
#### Swagger API Documentation
For convenient testing of the API’s functionality, a **Swagger interface** is available [here](http://localhost:8080/swagger/index.html#/). Swagger provides an interactive documentation interface where you can explore, test, and view the responses of each endpoint without needing to manually set headers or construct requests.

#### WebSocket Testing
To test WebSocket functionality that requires an authorization token, you can use **Postman** or another WebSocket client that supports adding headers.

##### 1. Setting Up the WebSocket Connection:
- In Postman, select **New WebSocket Request**.
- Enter the WebSocket URL:
    ```bash
    ws://localhost:8080/ws?group_id={id-of-group}
- Replace `{id-of-group}` with the group ID you wish to connect to.
##### 2. Adding the Authorization Header:

- In the **Headers** section, add:
    - **Key**: `Authorization`
    - **Value**: `Bearer {your-token}`
- Replace `{your-token}` with the actual token you received from the server.
##### 3. Connecting and Testing:

- Click **Connect** to start the WebSocket session.
- Upon successful connection, you can send and receive messages.
- If an error occurs, ensure your token is valid and the `group_id` is correct.
