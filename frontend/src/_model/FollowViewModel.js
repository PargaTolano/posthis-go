export class FollowViewModel {
    followerID;
    followedID;

    constructor({followerID, followedID}){
        this.followerID = followerID;
        this.followedID = followedID;
    }
}

export default FollowViewModel;