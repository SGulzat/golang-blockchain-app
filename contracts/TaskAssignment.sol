pragma solidity ^0.8.0;

contract TaskAssignment {
    struct Developer {
        uint256 id;  // ID разработчика из MySQL
        string level;  // "Junior", "Middle", "Senior"
    }

    struct Task {
        uint256 taskId;
        string taskLevel; // "Easy", "Medium", "Hard"
    }

    Developer[] public developers;
    mapping(uint256 => uint256) public taskAssignments; // taskId => developerId

    // Функция для добавления разработчика
    function addDeveloper(uint256 _id, string memory _level) public {
        developers.push(Developer(_id, _level));
    }

    // Логика назначения задачи разработчику
    function assignTask(uint256 _taskId, string memory _taskLevel) public returns (uint256) {
        for (uint256 i = 0; i < developers.length; i++) {
            if (keccak256(bytes(_taskLevel)) == keccak256(bytes("Easy")) && keccak256(bytes(developers[i].level)) == keccak256(bytes("Junior"))) {
                taskAssignments[_taskId] = developers[i].id;
                return developers[i].id;
            } else if (keccak256(bytes(_taskLevel)) == keccak256(bytes("Medium")) && keccak256(bytes(developers[i].level)) == keccak256(bytes("Middle"))) {
                taskAssignments[_taskId] = developers[i].id;
                return developers[i].id;
            } else if (keccak256(bytes(_taskLevel)) == keccak256(bytes("Hard")) && keccak256(bytes(developers[i].level)) == keccak256(bytes("Senior"))) {
                taskAssignments[_taskId] = developers[i].id;
                return developers[i].id;
            }
        }
        return 0; // Если никто не подходит
    }

    // Получение ID назначенного разработчика
    function getAssignedDeveloper(uint256 _taskId) public view returns (uint256) {
        return taskAssignments[_taskId];
    }
}