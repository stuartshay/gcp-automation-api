#!/bin/bash

# GitHub Cache Cleanup Script
# This script helps manage GitHub Actions cache storage when approaching limits

set -e

# Configuration
REPO="${GITHUB_REPOSITORY:-${1:-stuartshay/gcp-automation-api}}"
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

# Helper function to get cutoff date with cross-platform compatibility
get_cutoff_date() {
    local days=$1
    # Try different date command formats for cross-platform compatibility
    date -d "${days} days ago" --iso-8601=seconds 2>/dev/null || \
    date -v-${days}d '+%Y-%m-%dT%H:%M:%S%z' 2>/dev/null || \
    date -d "${days} days ago" '+%Y-%m-%dT%H:%M:%S%z' 2>/dev/null || {
        log_error "Unable to calculate cutoff date. Please ensure 'date' command supports ISO 8601 format."
        exit 1
    }
}

# Helper function to delete cache with proper error handling
delete_cache() {
    local cache_id=$1
    local description=$2

    if gh api --method DELETE repos/$REPO/actions/caches/$cache_id 2>/dev/null; then
        echo "  ✓ Deleted: $description"
        return 0
    else
        local exit_code=$?
        if [ $exit_code -eq 22 ]; then
            # HTTP 404 - cache already deleted
            echo "  ⚠ Cache already deleted: $description"
            return 0
        else
            log_warning "Failed to delete cache (exit code: $exit_code): $description"
            return 1
        fi
    fi
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

    local cutoff_date
    cutoff_date=$(get_cutoff_date 7)

    local deleted_count=0
    local failed_count=0

    # Use process substitution to avoid subshell issue with counter
    while IFS=' ' read -r cache_id key created_at; do
        if [ -n "$cache_id" ]; then
            if delete_cache "$cache_id" "$key (created: $created_at)"; then
                ((deleted_count++))
            else
                ((failed_count++))
            fi
        fi
    done < <(echo "$CACHES" | jq -r --arg cutoff "$cutoff_date" \
        '.actions_caches[] | select(.created_at < $cutoff) | "\(.id) \(.key) \(.created_at)"')

    log_success "Conservative cleanup completed: $deleted_count deleted, $failed_count failed"
}

# Aggressive cleanup (3+ days old, then oldest first until under limit)
cleanup_aggressive() {
    log_info "Running aggressive cleanup..."

    local max_size_bytes
    max_size_bytes=$(echo "$MAX_SIZE_GB * 1024 * 1024 * 1024" | bc)
    local cutoff_date
    cutoff_date=$(get_cutoff_date 3)

    local deleted_count=0
    local failed_count=0

    # Delete old caches first
    log_info "Deleting caches older than 3 days..."
    while IFS=' ' read -r cache_id key; do
        if [ -n "$cache_id" ]; then
            if delete_cache "$cache_id" "$key (old cache)"; then
                ((deleted_count++))
            else
                ((failed_count++))
            fi
        fi
    done < <(echo "$CACHES" | jq -r --arg cutoff "$cutoff_date" \
        '.actions_caches[] | select(.created_at < $cutoff) | "\(.id) \(.key)"')

    # Check current size and delete more if needed
    local current_caches
    current_caches=$(gh api repos/$REPO/actions/caches --paginate)
    local current_size
    current_size=$(echo "$current_caches" | jq '[.actions_caches[].size_in_bytes] | add // 0')

    if [ "$current_size" -gt "$max_size_bytes" ]; then
        log_info "Still over limit, deleting more caches (oldest first)..."

        # Calculate remaining size needed to delete
        local remaining_size_to_delete
        remaining_size_to_delete=$(echo "$current_size - $max_size_bytes" | bc)
        local running_size=0

        while IFS=' ' read -r created_at cache_id key size; do
            if [ "$running_size" -lt "$remaining_size_to_delete" ]; then
                if delete_cache "$cache_id" "$key (created: $created_at, size: $size bytes)"; then
                    ((deleted_count++))
                    running_size=$(echo "$running_size + $size" | bc)
                else
                    ((failed_count++))
                fi
            else
                log_success "Target size reached, stopping cleanup"
                break
            fi
        done < <(echo "$current_caches" | jq -r '.actions_caches[] | "\(.created_at) \(.id) \(.key) \(.size_in_bytes)"' | sort)
    fi

    log_success "Aggressive cleanup completed: $deleted_count deleted, $failed_count failed"
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
    local deleted_count=0
    local failed_count=0

    while IFS=' ' read -r cache_id key; do
        if [ -n "$cache_id" ]; then
            if delete_cache "$cache_id" "$key"; then
                ((deleted_count++))
            else
                ((failed_count++))
            fi
        fi
    done < <(echo "$CACHES" | jq -r '.actions_caches[] | "\(.id) \(.key)"')

    log_success "Emergency cleanup completed: $deleted_count deleted, $failed_count failed - all caches deleted"
}

# Clean specific cache types
cleanup_specific_types() {
    log_info "Cleaning up specific problematic cache types..."

    local deleted_count=0
    local failed_count=0

    # Delete Docker buildx caches (usually large)
    echo "Deleting Docker buildx caches..."
    while IFS=' ' read -r cache_id key; do
        if [ -n "$cache_id" ]; then
            if delete_cache "$cache_id" "buildx cache: $key"; then
                ((deleted_count++))
            else
                ((failed_count++))
            fi
        fi
    done < <(echo "$CACHES" | jq -r '.actions_caches[] | select(.key | contains("buildx")) | "\(.id) \(.key)"' | head -20)

    # Clean up duplicate Go caches
    echo "Cleaning up duplicate Go module caches..."
    while read -r cache_id; do
        if [ -n "$cache_id" ]; then
            if delete_cache "$cache_id" "duplicate Go cache: $cache_id"; then
                ((deleted_count++))
            else
                ((failed_count++))
            fi
        fi
    done < <(echo "$CACHES" | jq -r '.actions_caches[] | select(.key | contains("go-")) | "\(.created_at) \(.key) \(.id)"' | \
        sort -r | \
        awk '{
            # Extract the hash part from the cache key
            if (match($2, /go-.*-([a-f0-9]+)$/, arr)) {
                hash = arr[1]
                if (seen[hash]++) {
                    print $3  # Print cache ID for deletion
                }
            }
        }')

    log_success "Specific cleanup completed: $deleted_count deleted, $failed_count failed"
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
