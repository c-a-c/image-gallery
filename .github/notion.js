// .github/notion.js

/**
 * Notion APIを使用してGitHub IssueをNotionに同期するスクリプト
 *
 * このスクリプトはGitHub Actionsで実行され、Issueが作成、更新されたときにNotionページを自動的に作成または更新します。
 *
 * スクリプトの動作:
 * 1. GitHub Actionsから渡されたIssue情報を取得
 * 2. Notionデータベースを検索し、Issue IDに対応するページを取得
 * 3. Issueが新規作成された場合、Notionページを作成
 * 4. Issueが更新された場合、Notionページを更新
 */
import { Client } from "@notionhq/client";

const NOTION_TOKEN = process.env.NOTION_TOKEN;
const NOTION_DATABASE_ID = process.env.NOTION_DATABASE_ID;

// GitHub Actionsから渡されたIssue情報を取得
const ISSUE_ACTION = process.env.ISSUE_ACTION;
const ISSUE_NUMBER = process.env.ISSUE_NUMBER;
const ISSUE_TITLE = process.env.ISSUE_TITLE || 'No Title';
const ISSUE_URL = process.env.ISSUE_URL || 'No URL.';
const ISSUE_STATE = process.env.ISSUE_STATE;
const ISSUE_LABELS_JSON = process.env.ISSUE_LABELS;

const notion = new Client({ auth: NOTION_TOKEN });

/**
 * GitHubのIssueステータスに基づき、Notionのステータス名を決定する
 * 🚨 ご自身のNotionデータベースのステータス名に合わせて変更してください
 */
function getNotionStatus(issueState) {
  if (issueState === 'closed') {
    return '完了'; // Issueがクローズされたら「完了」
  }
  // その他の場合は「対応中」または「未対応」に自動設定
  return '対応中';
}

/**
 * Issue IDでNotionデータベースを検索し、対応するページIDを取得する
 */
async function findNotionPage(issueNumber) {
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
    return response.results[0].id;
  }
  return null;
}

/**
 * Notionページを作成または更新するための共通プロパティを構築する
 */
function buildNotionProperties(isNew = false) {
  const notionStatus = getNotionStatus(ISSUE_STATE);

  const properties = {
    // 1. タイトルプロパティ (必須)
    "Name": {
      title: [{ text: { content: ISSUE_TITLE } }],
    },
    // 2. ステータスプロパティ (Select/Statusタイプ)
    "ステータス": {
      status: {
        name: notionStatus,
      },
    },
    // 3. GitHub URLプロパティ (URLタイプ)
    "GitHub URL": {
      url: ISSUE_URL
    },
    // 4. ラベルプロパティ (Multi-Selectタイプ)
    // ※DBに「ラベル」プロパティ(Multi-Select)がある場合
    "ラベル": {
        multi_select: JSON.parse(ISSUE_LABELS_JSON).map(name => ({ name })),
    },
  };

  if (isNew) {
      // 新規作成時のみ、Issue IDプロパティを設定
      properties["Issue ID"] = { // 🚨 Notion側のIssue IDプロパティ名に合わせる
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
    const properties = buildNotionProperties(true);

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
    const properties = buildNotionProperties(false);

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
 * メイン処理
 */
async function main() {
  const pageId = await findNotionPage(ISSUE_NUMBER);

  if (ISSUE_ACTION === 'opened') {
    // Issueが新規作成された場合、Notionページが存在しないことを確認してから作成
    if (pageId) {
        console.log("警告: Issue IDに対応するページが既に存在します。スキップします。");
        return;
    }
    await createNewNotionPage();
  }
  else if (pageId) {
    // 既存のIssueが更新された場合（edited, closed, labeledなど）、ページを更新
    await updateNotionPage(pageId);
  } else {
    // 更新アクションだが、対応するNotionページが見つからない場合
    console.log(`警告: アクションは ${ISSUE_ACTION} ですが、Issue ID ${ISSUE_NUMBER} に対応するNotionページが見つかりませんでした。`);
  }
}

main();
