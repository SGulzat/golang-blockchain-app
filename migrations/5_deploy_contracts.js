const TaskAssignmentPro = artifacts.require("TaskAssignmentPro");

module.exports = function(deployer) {
    deployer.deploy(TaskAssignmentPro);
};