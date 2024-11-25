pragma solidity ^0.8.0;

contract TaskAssignmentPro {
    enum DeveloperLevel { Junior, Middle, Senior }
    enum TaskLevel { Easy, Medium, Hard }

    struct Developer {
        uint256 id;
        DeveloperLevel level;
    }

    mapping(uint256 => Developer) public developers;
    uint256 public developerCount;

    mapping(uint256 => uint256) public taskAssignments;

    event TaskAssigned(uint256 taskId, uint256 developerId);
    event Debug(string message);

    function addDeveloper(uint256 _id, DeveloperLevel _level) public {
        developers[developerCount] = Developer(_id, _level);
        developerCount++;
    }

    function assignTask(uint256 _taskId, TaskLevel _taskLevel) public returns (uint256) {
        for (uint256 i = 0; i < developerCount; i++) {
            if (developers[i].level == _getDeveloperLevelForTask(_taskLevel)) {
                taskAssignments[_taskId] = developers[i].id;
                emit TaskAssigned(_taskId, developers[i].id);
                return developers[i].id;
            }
        }
        emit Debug("No developer matches the task level");
        revert("No developer matches the task level");
    }

    function getAssignedDeveloper(uint256 _taskId) public view returns (uint256) {
        return taskAssignments[_taskId];
    }

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
}