// .github/notion.js

/**
 * Notion APIを使用してGitHub IssueをNotionに同期するスクリプト
 *
 * スクリプトの動作:
 * 1. GitHub Actionsから渡されたIssue情報を取得
 * 2. Notionデータベースを検索し、Issue IDに対応するページを取得
 * 3. Issueが新規作成された場合、Notionページを作成
 * 4. Issueが更新された場合、Notionページを更新
 */

import { Client } from "@notionhq/client";

// 環境変数の取得とバリデーション
const NOTION_TOKEN = process.env.NOTION_TOKEN;
const NOTION_DATABASE_ID = process.env.NOTION_DATABASE_ID;

// GitHub Actionsから渡されたIssue情報を取得
const ISSUE_ACTION = process.env.ISSUE_ACTION;
const ISSUE_NUMBER = process.env.ISSUE_NUMBER;
const ISSUE_TITLE = process.env.ISSUE_TITLE || 'No Title';
const ISSUE_URL = process.env.ISSUE_URL || 'No URL.';
const ISSUE_STATE = process.env.ISSUE_STATE;
const ISSUE_ASSIGNEES_JSON = process.env.ISSUE_ASSIGNEES;

// 環境変数の検証
if (!NOTION_TOKEN) {
  console.error("❌ エラー: NOTION_TOKEN が設定されていません。");
  process.exit(1);
}

if (!NOTION_DATABASE_ID) {
  console.error("❌ エラー: NOTION_DATABASE_ID が設定されていません。");
  process.exit(1);
}

console.log("🔧 Notion クライアントを初期化中...");
console.log(`🎬 アクション: ${ISSUE_ACTION}`);
console.log(`🔢 Issue番号: ${ISSUE_NUMBER}`);

let issueAssignees = [];
try {
  if (ISSUE_ASSIGNEES_JSON) {
    const parsed = JSON.parse(ISSUE_ASSIGNEES_JSON);
    issueAssignees = Array.isArray(parsed) ? parsed.map((name) => String(name)) : [];
  }
} catch (error) {
  console.warn("⚠️ Assignee情報の解析に失敗しました。空の配列として扱います:", error);
  issueAssignees = [];
}

const notion = new Client({ auth: NOTION_TOKEN });

/**
 * GitHubのIssueステータスに基づき、Notionのステータス名を決定する
 */
function getNotionStatus(issueState) {
  if (issueState === 'open') {
    return 'In progress';
  }
  if (issueState === 'closed') {
    return 'Done';
  }
  return 'Not started';
}

/**
 * Issue IDでNotionデータベースを検索し、対応するページIDを取得する
 */
async function findNotionPage(issueNumber) {
  try {
    console.log(`🔍 Issue ID: ${issueNumber} を検索中...`);

    const response = await notion.databases.query({
      database_id: NOTION_DATABASE_ID,
      filter: {
        property: "Issue ID",
        rich_text: {
          equals: String(issueNumber),
        },
      },
    });

    // 最初の結果のページIDを返す
    if (response.results.length > 0) {
      console.log(`✅ 既存ページが見つかりました (Page ID: ${response.results[0].id.substring(0, 8)}...)`);
      return response.results[0].id;
    }

    console.log(`Issue ID: ${issueNumber} に対応するページは見つかりませんでした`);
    return null;
  } catch (error) {
    console.error("❌ Notionデータベースの検索中にエラーが発生しました:");
    console.error("エラー詳細:", error.message);
    if (error.body) {
      console.error("APIレスポンス:", JSON.stringify(error.body, null, 2));
    }
    throw error;
  }
}

/**
 * Notionページを作成または更新するための共通プロパティを構築する
 */
function buildNotionProperties(options = {}) {
  const { includeIssueId = false } = options;
  const notionStatus = getNotionStatus(ISSUE_STATE);

  const properties = {
    "Title": {
      title: [{ text: { content: ISSUE_TITLE } }],
    },
    "Status": {
      status: {
        name: notionStatus,
      },
    },
    "URL": {
      url: ISSUE_URL,
    },
    "Assignee": {
      multi_select: issueAssignees.map((name) => ({ name })),
    },
  };

  if (includeIssueId) {
    properties["Issue ID"] = {
      rich_text: [{ text: { content: String(ISSUE_NUMBER) } }],
    };
  }

  return properties;
}

/**
 * 新しいNotionページを作成する
 */
async function createNewNotionPage() {
  console.log(`アクション: ${ISSUE_ACTION} - 新規ページを作成します`);
  try {
    const properties = buildNotionProperties({ includeIssueId: true });

    const response = await notion.pages.create({
      parent: { database_id: NOTION_DATABASE_ID },
      properties: properties,
      // Issue本文はここでは省略 (更新時に本文全体を上書きするのは複雑なため)
    });
    console.log("✅ Notionページが正常に作成されました。");
  } catch (error) {
    console.error("❌ Notionページの作成に失敗しました:", error.body || error);
  }
}

/**
 * 既存のNotionページを更新する
 */
async function updateNotionPage(pageId) {
  console.log(`アクション: ${ISSUE_ACTION} - 既存ページを更新します (Page ID: ${pageId})`);
  try {
    const properties = buildNotionProperties();

    const response = await notion.pages.update({
      page_id: pageId,
      properties: properties,
    });
    console.log("✅ Notionページが正常に更新されました。");
  } catch (error) {
    console.error("❌ Notionページの更新に失敗しました:", error.body || error);
  }
}

/**
 * 既存のNotionページのステータスのみを更新する
 */
async function updateNotionStatus(pageId, issueState) {
  console.log(`アクション: ${ISSUE_ACTION} - ステータスのみ更新します (Page ID: ${pageId})`);
  try {
    const statusName = getNotionStatus(issueState);
    await notion.pages.update({
      page_id: pageId,
      properties: {
        "Status": {
          status: { name: statusName },
        },
      },
    });
    console.log("✅ ステータスのみ更新しました。");
  } catch (error) {
    console.error("❌ ステータス更新に失敗しました:", error.body || error);
  }
}

/**
 * メイン処理
 *
 * 動作フロー:
 * 1. Issue IDでNotionデータベース内の既存ページを検索
 * 2. Issueが新規作成（opened）の場合:
 *    - ページが存在しない → 新規作成
 *    - ページが既に存在 → スキップ（警告表示）
 * 3. Issueが更新（edited, closed, reopened）の場合:
 *    - ページが存在する → 既存ページを更新（タイトル、ステータス、Assigneeなどを更新）
 *    - ページが存在しない → 新規作成
 */
async function main() {
  try {
    console.log("\n" + "=".repeat(60));
    console.log("🚀 GitHub Issue → Notion 同期処理を開始");
    console.log("=".repeat(60));

    const pageId = await findNotionPage(ISSUE_NUMBER);

    const isOpened = ISSUE_ACTION === 'opened';
    const isClosed = ISSUE_ACTION === 'closed';

    if (isOpened) {
      // Issueが新規作成された場合
      if (pageId) {
        console.log("⚠️ 警告: Issue IDに対応するページが既に存在します。スキップします。");
        return;
      }
      await createNewNotionPage();
    } else if (!pageId) {
      console.log(`Issue ID ${ISSUE_NUMBER} に対応するNotionページが見つかりませんでした。`);
      if (isClosed) {
        console.log("⚠️ クローズされたIssueですが対応するタスクが存在しないため、処理をスキップします。");
        return;
      }
      console.log(`新規ページを作成します...`);
      await createNewNotionPage();
    } else if (isClosed) {
      await updateNotionStatus(pageId, ISSUE_STATE);
    } else {
      console.log(`既存のNotionページを更新します...`);
      await updateNotionPage(pageId);
    }

    console.log("\n" + "=".repeat(60));
    console.log("処理が正常に完了しました");
    console.log("=".repeat(60) + "\n");

  } catch (error) {
    console.error("\n" + "=".repeat(60));
    console.error("処理中に致命的なエラーが発生しました");
    console.error("=".repeat(60));
    console.error("エラー:", error);
    process.exit(1);
  }
}

main();
