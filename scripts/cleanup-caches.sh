#!/bin/bash

# GitHub Cache Cleanup Script
# This script helps manage GitHub Actions cache storage when approaching limits

set -e

# Configuration
REPO="${1:-stuartshay/gcp-automation-api}"
MAX_SIZE_GB="${2:-8}"
DEFAULT_STRATEGY="${3:-conservative}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if gh CLI is installed
check_dependencies() {
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI (gh) is not installed. Please install it first:"
        echo "  brew install gh  # macOS"
        echo "  sudo apt install gh  # Ubuntu/Debian"
        echo "  winget install GitHub.cli  # Windows"
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        log_error "jq is not installed. Please install it first:"
        echo "  brew install jq  # macOS"
        echo "  sudo apt install jq  # Ubuntu/Debian"
        echo "  winget install jqlang.jq  # Windows"
        exit 1
    fi

    if ! command -v bc &> /dev/null; then
        log_error "bc is not installed. Please install it first:"
        echo "  brew install bc  # macOS"
        echo "  sudo apt install bc  # Ubuntu/Debian"
        exit 1
    fi
}

# Check authentication
check_auth() {
    if ! gh auth status &> /dev/null; then
        log_error "Not authenticated with GitHub CLI. Please run:"
        echo "  gh auth login"
        exit 1
    fi
}

# Get current cache information
get_cache_info() {
    log_info "Fetching cache information for repository: $REPO"

    CACHES=$(gh api repos/$REPO/actions/caches --paginate 2>/dev/null)
    if [ $? -ne 0 ]; then
        log_error "Failed to fetch cache information. Check repository name and permissions."
        exit 1
    fi

    TOTAL_SIZE_BYTES=$(echo "$CACHES" | jq '[.actions_caches[].size_in_bytes] | add // 0')
    TOTAL_SIZE_GB=$(echo "scale=2; $TOTAL_SIZE_BYTES / 1024 / 1024 / 1024" | bc)
    CACHE_COUNT=$(echo "$CACHES" | jq '.actions_caches | length')

    echo "Current cache status:"
    echo "  Total size: ${TOTAL_SIZE_GB} GB"
    echo "  Cache count: ${CACHE_COUNT}"
    echo "  Limit: 10 GB"

    if (( $(echo "$TOTAL_SIZE_GB > 9" | bc -l) )); then
        log_warning "Cache usage is critically high!"
    elif (( $(echo "$TOTAL_SIZE_GB > 7" | bc -l) )); then
        log_warning "Cache usage is approaching limit"
    else
        log_success "Cache usage is within acceptable limits"
    fi
}

# List caches by type and size
list_caches() {
    log_info "Listing caches by type and size:"
    echo
    echo "Top 20 largest caches:"
    echo "$CACHES" | jq -r '.actions_caches[] | "\((.size_in_bytes / 1024 / 1024 / 1024 * 100 | floor) / 100) GB | \(.key) | \(.created_at)"' | \
    sort -rn | head -20 | \
    while read -r size unit key created; do
        printf "  %6s %s | %s | %s\n" "$size" "$unit" "$key" "$created"
    done
    echo

    echo "Caches by type:"
    echo "$CACHES" | jq -r '.actions_caches[].key' | \
    sed 's/-[a-f0-9]*$//' | sort | uniq -c | sort -rn | \
    while read -r count type; do
        printf "  %3d caches | %s\n" "$count" "$type"
    done
    echo
}

# Conservative cleanup (7+ days old)
cleanup_conservative() {
    log_info "Running conservative cleanup (deleting caches older than 7 days)..."

    CUTOFF_DATE=$(date -d '7 days ago' --iso-8601=seconds 2>/dev/null || date -v-7d '+%Y-%m-%dT%H:%M:%S%z' 2>/dev/null || date -d '7 days ago' '+%Y-%m-%dT%H:%M:%S%z')

    DELETED_COUNT=0
    echo "$CACHES" | jq -r --arg cutoff "$CUTOFF_DATE" \
    '.actions_caches[] | select(.created_at < $cutoff) | "\(.id) \(.key) \(.created_at)"' | \
    while read -r cache_id key created_at; do
        if [ -n "$cache_id" ]; then
            echo "  Deleting: $key (created: $created_at)"
            if gh api --method DELETE repos/$REPO/actions/caches/$cache_id &> /dev/null; then
                ((DELETED_COUNT++))
            else
                log_warning "Failed to delete cache: $cache_id"
            fi
        fi
    done

    log_success "Conservative cleanup completed"
}

# Aggressive cleanup (3+ days old, then oldest first until under limit)
cleanup_aggressive() {
    log_info "Running aggressive cleanup..."

    MAX_SIZE_BYTES=$(echo "$MAX_SIZE_GB * 1024 * 1024 * 1024" | bc)
    CUTOFF_DATE=$(date -d '3 days ago' --iso-8601=seconds 2>/dev/null || date -v-3d '+%Y-%m-%dT%H:%M:%S%z' 2>/dev/null || date -d '3 days ago' '+%Y-%m-%dT%H:%M:%S%z')

    # Delete old caches first
    log_info "Deleting caches older than 3 days..."
    echo "$CACHES" | jq -r --arg cutoff "$CUTOFF_DATE" \
    '.actions_caches[] | select(.created_at < $cutoff) | "\(.id) \(.key)"' | \
    while read -r cache_id key; do
        if [ -n "$cache_id" ]; then
            echo "  Deleting old cache: $key"
            gh api --method DELETE repos/$REPO/actions/caches/$cache_id &> /dev/null || true
        fi
    done

    # Check current size and delete more if needed
    CURRENT_CACHES=$(gh api repos/$REPO/actions/caches --paginate)
    CURRENT_SIZE=$(echo "$CURRENT_CACHES" | jq '[.actions_caches[].size_in_bytes] | add // 0')

    if [ "$CURRENT_SIZE" -gt "$MAX_SIZE_BYTES" ]; then
        log_info "Still over limit, deleting more caches (oldest first)..."
        echo "$CURRENT_CACHES" | jq -r '.actions_caches[] | "\(.created_at) \(.id) \(.key)"' | \
        sort | \
        while read -r created_at cache_id key; do
            CURRENT_SIZE=$(gh api repos/$REPO/actions/caches --paginate | jq '[.actions_caches[].size_in_bytes] | add // 0')
            if [ "$CURRENT_SIZE" -gt "$MAX_SIZE_BYTES" ]; then
                echo "  Deleting: $key (created: $created_at)"
                gh api --method DELETE repos/$REPO/actions/caches/$cache_id &> /dev/null || true
            else
                log_success "Target size reached"
                break
            fi
        done
    fi

    log_success "Aggressive cleanup completed"
}

# Emergency cleanup (delete all caches)
cleanup_emergency() {
    log_warning "Running EMERGENCY cleanup - this will delete ALL caches!"
    read -p "Are you sure you want to delete ALL caches? This cannot be undone. [y/N]: " -n 1 -r
    echo

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Emergency cleanup cancelled"
        return
    fi

    log_info "Deleting all caches..."
    echo "$CACHES" | jq -r '.actions_caches[] | "\(.id) \(.key)"' | \
    while read -r cache_id key; do
        if [ -n "$cache_id" ]; then
            echo "  Deleting: $key"
            gh api --method DELETE repos/$REPO/actions/caches/$cache_id &> /dev/null || true
        fi
    done

    log_success "Emergency cleanup completed - all caches deleted"
}

# Clean specific cache types
cleanup_specific_types() {
    log_info "Cleaning up specific problematic cache types..."

    # Delete Docker buildx caches (usually large)
    echo "Deleting Docker buildx caches..."
    echo "$CACHES" | jq -r '.actions_caches[] | select(.key | contains("buildx")) | "\(.id) \(.key)"' | \
    head -20 | \
    while read -r cache_id key; do
        if [ -n "$cache_id" ]; then
            echo "  Deleting buildx cache: $key"
            gh api --method DELETE repos/$REPO/actions/caches/$cache_id &> /dev/null || true
        fi
    done

    # Clean up duplicate Go caches
    echo "Cleaning up duplicate Go module caches..."
    echo "$CACHES" | jq -r '.actions_caches[] | select(.key | contains("go-")) | "\(.created_at) \(.key) \(.id)"' | \
    sort -r | \
    awk '{
        # Extract the hash part from the cache key
        if (match($2, /go-.*-([a-f0-9]+)$/, arr)) {
            hash = arr[1]
            if (seen[hash]++) {
                print $3  # Print cache ID for deletion
            }
        }
    }' | \
    while read -r cache_id; do
        if [ -n "$cache_id" ]; then
            echo "  Deleting duplicate Go cache: $cache_id"
            gh api --method DELETE repos/$REPO/actions/caches/$cache_id &> /dev/null || true
        fi
    done
}

# Main function
main() {
    echo "GitHub Actions Cache Cleanup Tool"
    echo "================================="
    echo

    check_dependencies
    check_auth
    get_cache_info

    if [ "$CACHE_COUNT" -eq 0 ]; then
        log_info "No caches found to clean up"
        exit 0
    fi

    echo
    echo "Cleanup strategies:"
    echo "  1. Conservative - Delete caches older than 7 days"
    echo "  2. Aggressive   - Delete caches older than 3 days, then oldest until under ${MAX_SIZE_GB}GB"
    echo "  3. Emergency    - Delete ALL caches (cannot be undone)"
    echo "  4. Specific     - Clean up problematic cache types only"
    echo "  5. List only    - Just show cache information"
    echo

    if [ "$DEFAULT_STRATEGY" != "conservative" ]; then
        STRATEGY="$DEFAULT_STRATEGY"
    else
        read -p "Choose cleanup strategy [1-5] (default: 1): " STRATEGY
        STRATEGY=${STRATEGY:-1}
    fi

    case $STRATEGY in
        1|conservative)
            list_caches
            cleanup_conservative
            ;;
        2|aggressive)
            list_caches
            cleanup_aggressive
            ;;
        3|emergency)
            list_caches
            cleanup_emergency
            ;;
        4|specific)
            list_caches
            cleanup_specific_types
            ;;
        5|list)
            list_caches
            exit 0
            ;;
        *)
            log_error "Invalid strategy selected"
            exit 1
            ;;
    esac

    echo
    log_info "Final cache status:"
    get_cache_info
}

# Show usage
show_usage() {
    echo "Usage: $0 [REPO] [MAX_SIZE_GB] [STRATEGY]"
    echo
    echo "Arguments:"
    echo "  REPO          Repository name (default: stuartshay/gcp-automation-api)"
    echo "  MAX_SIZE_GB   Maximum cache size in GB (default: 8)"
    echo "  STRATEGY      Cleanup strategy: conservative|aggressive|emergency|specific|list"
    echo
    echo "Examples:"
    echo "  $0                                           # Interactive mode"
    echo "  $0 owner/repo 7 aggressive                   # Aggressive cleanup to 7GB"
    echo "  $0 owner/repo 8 list                         # List caches only"
    echo
    echo "Environment variables:"
    echo "  GITHUB_TOKEN  GitHub token (if not using gh auth)"
    echo
}

# Handle command line arguments
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    show_usage
    exit 0
fi

# Run main function
main "$@"
