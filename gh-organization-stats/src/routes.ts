import express from "express";
import axios from "axios";
import ClubStats from "./model";

const router = express.Router();

// Fetch # of commits for a repository since the given date
const fetchCommits = async (owner: string, repo: string, date: string) => {
  try {
    const response = await axios.get(
      `https://api.github.com/repos/${owner}/${repo}/commits`,
      {
        params: { since: date },
        headers: { Accept: "application/vnd.github.v3+json" },
      },
    );

    return response.data.length;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const statusCode = error.response?.status;

      if (statusCode === 404) {
        throw { message: `Repository ${owner}/${repo} not found.`, statusCode };
      }

      throw {
        message: `Error fetching commits for ${owner}/${repo}.`,
        statusCode,
      };
    }
    console.error(`Error fetching commits for ${owner}/${repo}:`, error);
    throw error;
  }
};

// Fetch # of opened pull requests for a repository since the given date
const fetchOpenedPullRequests = async (
  owner: string,
  repo: string,
  date: string,
) => {
  try {
    const response = await axios.get(
      `https://api.github.com/repos/${owner}/${repo}/pulls`,
      {
        params: { state: "all" },
        headers: { Accept: "application/vnd.github.v3+json" },
      },
    );

    let openedPRs = 0;
    for (const pr of response.data) {
      const createdAt = new Date(pr.created_at);
      if (createdAt >= new Date(date)) {
        openedPRs++;
      }
    }

    return openedPRs;
  } catch (error) {
    console.error(`Error fetching PRs for ${owner}/${repo}:`, error);
    return 0;
  }
};

// Route to get stats and save club-wide totals to the database
router.post("/stats", async (req, res) => {
  const { startDate, repos } = req.body;

  if (!startDate || !repos) {
    return res
      .status(400)
      .json({ error: "Please provide startDate and repos." });
  }

  let totalCommits = 0;
  let totalOpenedPRs = 0;

  for (const [owner, repo] of repos) {
    try {
      totalCommits += await fetchCommits(owner, repo, startDate as string);
      totalOpenedPRs += await fetchOpenedPullRequests(
        owner,
        repo,
        startDate as string,
      );
    } catch (error) {
      if (error && typeof error === "object" && "message" in error) {
        const err = error as { message: string; statusCode?: number };
        const statusCode = err.statusCode || 500;
        console.error(
          `Error fetching stats for ${owner}/${repo}: ${err.message}`,
        );
        return res.status(statusCode).json({
          error: `Error fetching stats for ${owner}/${repo}: ${err.message}`,
        });
      } else {
        console.error(`Unexpected error fetching stats for ${owner}/${repo}.`);
        return res.status(500).json({ error: "Unexpected error occurred." });
      }
    }
  }

  const today = new Date(new Date().setUTCHours(0, 0, 0, 0));

  try {
    let clubStat = await ClubStats.findOne({ start_date: startDate });

    if (clubStat) {
      clubStat.totalCommits = totalCommits;
      clubStat.totalOpenedPRs = totalOpenedPRs;
      clubStat.repos = repos;
      clubStat.collection_date = today;
      await clubStat.save();
    } else {
      clubStat = new ClubStats({
        collection_date: today,
        start_date: startDate,
        totalCommits,
        totalOpenedPRs,
        repos,
      });
      await clubStat.save();
    }

    res.json({
      message: "Club-wide stats recorded successfully",
      stats: {
        totalCommits: clubStat.totalCommits,
        totalOpenedPRs: clubStat.totalOpenedPRs,
      },
    });
  } catch (error) {
    console.error("Error saving club stats:", error);
    res.status(500).json({ error: "Error saving club stats" });
  }
});

router.get("/stats", async (req, res) => {
  try {
    const stats = await ClubStats.find().sort({ date: -1 }); // Fetch all records, sorted by date in descending order
    res.json({ stats });
  } catch (error) {
    console.error("Error fetching stats:", error);
    res.status(500).json({ error: "Error fetching stats" });
  }
});

export default router;
