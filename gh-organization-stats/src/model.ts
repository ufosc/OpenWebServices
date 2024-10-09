import mongoose from "mongoose";

// One record per semester (start date is primary key)
const ClubStatsSchema = new mongoose.Schema({
  start_date: Date,
  totalCommits: Number,
  totalOpenedPRs: Number,
  repos: Array,
  collection_date: {
    type: Date,
    default: () => new Date(new Date().setUTCHours(0, 0, 0, 0)),
  },
});

const ClubStats = mongoose.model("ClubStats", ClubStatsSchema);
export default ClubStats;
