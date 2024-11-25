const TaskAssignment = artifacts.require("TaskAssignment");

module.exports = function(deployer) {
    deployer.deploy(TaskAssignment);
};