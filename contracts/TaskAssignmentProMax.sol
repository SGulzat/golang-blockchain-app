pragma solidity ^0.8.0;

contract TaskAssignmentProMax {
    enum DeveloperLevel { Junior, Middle, Senior }
    enum TaskLevel { Easy, Medium, Hard }

    struct Developer {
        uint256 id;
        DeveloperLevel level;
    }

    // Хранение разработчиков по их уникальным ID
    mapping(uint256 => Developer) public developers;
    uint256[] public developerIds;

    mapping(uint256 => uint256) public taskAssignments;

    event TaskAssigned(uint256 taskId, uint256 developerId);
    event Debug(string message);

    // Добавление разработчика с уникальным ID
    function addDeveloper(uint256 _id, DeveloperLevel _level) public {
        // Проверка, что разработчик с таким ID еще не добавлен
        require(developers[_id].id == 0, "Developer already exists");

        developers[_id] = Developer(_id, _level);
        developerIds.push(_id);
    }

    // Назначение задачи разработчику соответствующего уровня
    function assignTask(uint256 _taskId, TaskLevel _taskLevel) public returns (uint256) {
        DeveloperLevel requiredLevel = _getDeveloperLevelForTask(_taskLevel);

        for (uint256 i = 0; i < developerIds.length; i++) {
            uint256 devId = developerIds[i];
            if (developers[devId].level == requiredLevel) {
                taskAssignments[_taskId] = devId;
                emit TaskAssigned(_taskId, devId);
                return devId;
            }
        }

        emit Debug("No developer matches the task level");
        revert("No developer matches the task level");
    }

    // Получение ID разработчика, назначенного на задачу
    function getAssignedDeveloper(uint256 _taskId) public view returns (uint256) {
        return taskAssignments[_taskId];
    }

    // Вспомогательная функция для определения уровня разработчика по уровню задачи
    function _getDeveloperLevelForTask(TaskLevel _taskLevel) internal pure returns (DeveloperLevel) {
        if (_taskLevel == TaskLevel.Easy) {
            return DeveloperLevel.Junior;
        } else if (_taskLevel == TaskLevel.Medium) {
            return DeveloperLevel.Middle;
        } else if (_taskLevel == TaskLevel.Hard) {
            return DeveloperLevel.Senior;
        } else {
            revert("Invalid task level");
        }
    }

    // (Опционально) Функция для сброса списка разработчиков
    function resetDevelopers() public {
        for (uint256 i = 0; i < developerIds.length; i++) {
            uint256 devId = developerIds[i];
            delete developers[devId];
        }
        delete developerIds;
    }
}