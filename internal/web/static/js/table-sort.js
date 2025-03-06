/**
 * TableSorter - A utility for making HTML tables sortable
 */
class TableSorter {
  constructor(tableId, options = {}) {
    this.table = document.getElementById(tableId);
    if (!this.table) {
      console.error(`Table with ID '${tableId}' not found`);
      return;
    }

    this.options = {
      headerClass: 'sortable',
      ascClass: 'sort-asc',
      descClass: 'sort-desc',
      initialColumn: null,
      initialDirection: 'asc',
      ...options
    };

    this.currentSort = {
      column: this.options.initialColumn,
      direction: this.options.initialDirection
    };

    this.init();
  }

  init() {
    // Get the headers from the table
    const thead = this.table.querySelector('thead');
    if (!thead) {
      console.error('Table does not have a thead element');
      return;
    }

    // Get all the headers
    const headers = thead.querySelectorAll('th');

    // Add click event listeners to headers
    headers.forEach((header, index) => {
      // Skip columns that have a 'data-sortable="false"' attribute
      if (header.getAttribute('data-sortable') === 'false') {
        return;
      }

      // Add the sortable class
      header.classList.add(this.options.headerClass);

      // Create sort indicator element
      const indicator = document.createElement('span');
      indicator.className = 'sort-indicator ml-2';
      indicator.innerHTML = '<i class="fas fa-sort text-gray-400"></i>';
      header.appendChild(indicator);

      // Set initial sort indicator if this is the initial sort column
      if (index === this.options.initialColumn) {
        this.updateSortIndicator(header, this.options.initialDirection);
      }

      // Add click event listener
      header.addEventListener('click', () => {
        this.sortTable(index, header);
      });
    });

    // Perform initial sort if specified
    if (this.options.initialColumn !== null) {
      this.sortTable(this.options.initialColumn, headers[this.options.initialColumn]);
    }
  }

  // Parse a value for sorting, converting to number if possible
  parseForSort(value) {
    // If it's already a number, return it
    if (typeof value === 'number') return value;

    // If it's a string that might represent a number
    if (typeof value === 'string') {
      // Remove commas and any other non-numeric chars except decimal point
      const cleanValue = value.replace(/[^\d.-]/g, '');
      // Try to convert to number
      const num = parseFloat(cleanValue);
      // Return the number if valid, otherwise the original string
      return isNaN(num) ? value.toLowerCase() : num;
    }

    // For other types, convert to string
    return String(value).toLowerCase();
  }

  sortTable(columnIndex, header) {
    const tbody = this.table.querySelector('tbody');
    if (!tbody) {
      console.error('Table does not have a tbody element');
      return;
    }

    // Determine sort direction
    let direction = 'asc';
    if (this.currentSort.column === columnIndex) {
      // Toggle direction if same column is clicked
      direction = this.currentSort.direction === 'asc' ? 'desc' : 'asc';
    }

    // Update current sort state
    this.currentSort.column = columnIndex;
    this.currentSort.direction = direction;

    // Update sort indicators on all headers
    const headers = this.table.querySelectorAll('th');
    headers.forEach(h => {
      // Reset all headers
      const indicator = h.querySelector('.sort-indicator');
      if (indicator) {
        indicator.innerHTML = '<i class="fas fa-sort text-gray-400"></i>';
      }
      h.classList.remove(this.options.ascClass, this.options.descClass);
    });

    // Update the clicked header
    this.updateSortIndicator(header, direction);

    // Get all rows for sorting
    const rows = Array.from(tbody.querySelectorAll('tr'));

    // Sort the rows
    rows.sort((rowA, rowB) => {
      // Get the cell values to compare
      const cellA = rowA.querySelectorAll('td')[columnIndex];
      const cellB = rowB.querySelectorAll('td')[columnIndex];

      if (!cellA || !cellB) return 0;

      // Check for data-sort-value attribute first
      const rawValueA = cellA.getAttribute('data-sort-value') || cellA.textContent.trim();
      const rawValueB = cellB.getAttribute('data-sort-value') || cellB.textContent.trim();

      // Parse values for sorting
      const valueA = this.parseForSort(rawValueA);
      const valueB = this.parseForSort(rawValueB);

      // Compare values
      let comparison = 0;
      if (typeof valueA === 'number' && typeof valueB === 'number') {
        // Numeric comparison
        comparison = valueA - valueB;
      } else {
        // String comparison
        comparison = String(valueA).localeCompare(String(valueB));
      }

      // Reverse for descending order
      return direction === 'asc' ? comparison : -comparison;
    });

    // Re-append rows in the sorted order
    rows.forEach(row => {
      tbody.appendChild(row);
    });
  }

  updateSortIndicator(header, direction) {
    const indicator = header.querySelector('.sort-indicator');
    if (indicator) {
      indicator.innerHTML = direction === 'asc'
        ? '<i class="fas fa-sort-up text-indigo-600"></i>'
        : '<i class="fas fa-sort-down text-indigo-600"></i>';
    }

    // Update classes
    header.classList.remove(this.options.ascClass, this.options.descClass);
    header.classList.add(direction === 'asc' ? this.options.ascClass : this.options.descClass);
  }
}

// Add CSS for sortable tables
document.addEventListener('DOMContentLoaded', () => {
  const style = document.createElement('style');
  style.textContent = `
    .sortable {
      cursor: pointer;
      user-select: none;
    }
    .sortable:hover {
      background-color: rgba(79, 70, 229, 0.1);
    }
    .sort-asc .sort-indicator,
    .sort-desc .sort-indicator {
      display: inline-block;
    }
  `;
  document.head.appendChild(style);
});
