use anyhow::{Context, Result};
use std::io::Write;
use std::process::{Command, Stdio, exit};

pub fn run_list(format: Option<String>, sort_by: Option<String>, starred: bool) -> Result<()> {
    let config = pado::GlobalConfig::load().unwrap_or_default();
    let list = pado::load_project_list().context("failed to load project list")?;

    let sort_order = sort_by.as_ref().unwrap_or(&config.defaults.sort_order);
    let output_format = format.as_ref().unwrap_or(&config.defaults.output_format);

    let projects = if starred {
        list.get_starred_projects()
    } else {
        list.get_projects_sorted_by(sort_order)
    };

    if projects.is_empty() {
        if starred {
            eprintln!("No starred projects. Use 'pd star' to star a project.");
        } else {
            eprintln!("No projects tracked yet. Use 'pd add' or navigate to a project with 'pd'.");
        }
        return Ok(());
    }

    match output_format.as_str() {
        "paths" => {
            for project in projects {
                println!("{}", project.path.display());
            }
        }
        "json" => {
            println!("{}", serde_json::to_string_pretty(&list)?);
        }
        _ => {
            use comfy_table::Table;

            let mut table = Table::new();
            table.set_header(vec!["Name", "Type", "Path", "Accessed", "Count", "★"]);

            for project in projects {
                let last_accessed = project.last_accessed.format("%Y-%m-%d %H:%M").to_string();
                let star = if project.starred { "★" } else { "" };
                table.add_row(vec![
                    &project.name,
                    &project.project_type,
                    &project.path.display().to_string(),
                    &last_accessed,
                    &project.access_count.to_string(),
                    star,
                ]);
            }

            println!("{}", table);
        }
    }
    Ok(())
}

pub fn run_recent(limit: Option<usize>) -> Result<()> {
    let config = pado::GlobalConfig::load().unwrap_or_default();
    let list = pado::load_project_list().context("failed to load project list")?;

    let limit = limit.unwrap_or(config.defaults.recent_limit);
    let projects = list.get_recent_projects(limit);

    if projects.is_empty() {
        eprintln!("No projects tracked yet.");
        return Ok(());
    }

    use comfy_table::Table;

    let mut table = Table::new();
    table.set_header(vec!["Name", "Type", "Path", "Last Accessed"]);

    for project in projects {
        let last_accessed = pado::format_time_ago(project.last_accessed.timestamp());
        table.add_row(vec![
            &project.name,
            &project.project_type,
            &project.path.display().to_string(),
            &last_accessed,
        ]);
    }

    println!("{}", table);
    Ok(())
}

pub fn run_stats() -> Result<()> {
    let list = pado::load_project_list().context("failed to load project list")?;

    let projects: Vec<_> = list.projects.values().collect();

    if projects.is_empty() {
        eprintln!("No projects tracked yet.");
        return Ok(());
    }

    println!("\nProject Statistics:\n");
    println!("Total projects: {}", projects.len());
    println!(
        "Starred projects: {}",
        projects.iter().filter(|p| p.starred).count()
    );

    let mut sorted_by_access: Vec<_> = projects.clone();
    sorted_by_access.sort_by(|a, b| b.access_count.cmp(&a.access_count));

    println!("\nMost Accessed:");
    for project in sorted_by_access.iter().take(10) {
        let star = if project.starred { "★ " } else { "  " };
        println!(
            "  {}{:<30} {} accesses",
            star, project.name, project.access_count
        );
    }

    let mut sorted_by_time: Vec<_> = projects.clone();
    sorted_by_time.sort_by(|a, b| b.last_accessed.cmp(&a.last_accessed));

    println!("\nRecently Accessed:");
    for project in sorted_by_time.iter().take(10) {
        let star = if project.starred { "★ " } else { "  " };
        let time_ago = pado::format_time_ago(project.last_accessed.timestamp());
        println!("  {}{:<30} {}", star, project.name, time_ago);
    }

    let mut type_counts: std::collections::HashMap<String, usize> =
        std::collections::HashMap::new();
    for project in &projects {
        *type_counts.entry(project.project_type.clone()).or_insert(0) += 1;
    }

    println!("\nProject Types:");
    let mut types: Vec<_> = type_counts.iter().collect();
    types.sort_by(|a, b| b.1.cmp(a.1));
    for (ptype, count) in types {
        println!("  {:<15} {}", ptype, count);
    }

    println!();
    Ok(())
}

pub fn run_switch(recent: bool, starred: bool) -> Result<()> {
    let config = pado::GlobalConfig::load().unwrap_or_default();
    let list = pado::load_project_list().context("failed to load project list")?;

    let projects = if starred {
        list.get_starred_projects()
    } else if recent {
        list.get_recent_projects(config.defaults.recent_limit)
    } else {
        list.get_projects_sorted_by(&config.defaults.sort_order)
    };

    if projects.is_empty() {
        if starred {
            eprintln!("No starred projects. Use 'pd star' to star a project.");
        } else if recent {
            eprintln!("No recent projects.");
        } else {
            eprintln!("No projects tracked yet. Use 'pd add' or navigate to a project with 'pd'.");
        }
        exit(1);
    }

    let mut fzf = Command::new("fzf")
        .arg("--height=40%")
        .arg("--reverse")
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .spawn()
        .context("failed to spawn fzf - is it installed?")?;

    if let Some(stdin) = fzf.stdin.as_mut() {
        for project in projects {
            let star = if project.starred { "★ " } else { "  " };
            writeln!(
                stdin,
                "{}{}\t{}\t{}",
                star,
                project.name,
                project.project_type,
                project.path.display()
            )?;
        }
    }

    let output = fzf.wait_with_output()?;

    if output.status.success() && !output.stdout.is_empty() {
        let selection = String::from_utf8_lossy(&output.stdout);
        if let Some(path) = selection.split('\t').nth(2) {
            let path = path.trim();
            println!("cd \"{}\"", path);
        }
    }
    Ok(())
}
