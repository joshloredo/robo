//! Path planning and navigation
use rppal::gpio::{Gpio, OutputPin};

use r2r::{Context, Node, Publisher, Subscriber, QosProfile};
use r2r::geometry_msgs::msg::{PoseStamped, Twist};
use r2r::nav_msgs::msg::Path;

//! Djikstra's algorithm

pub struct PathPlanner {
    context: Context,
    node: Node,
    publisher: Publisher<Path>,
    subscriber: Subscriber<Twist>,
}

impl PathPlanner { 
    pub fn new(context: Context, node: Node) -> Self {
        let publisher = node.create_publisher::<Path>("path", QosProfile::default());
        let subscriber = node.create_subscriber::<Twist>("cmd_vel", QosProfile::default());

        Self { context, node, publisher, subscriber }
    }

    pub fn plan_path(&self, start: PoseStamped, goal: PoseStamped) -> Result<Path, PathPlannerError> {
        fn heuristic(a: &PoseStamped, b: &PoseStamped) -> f32 {
            let dx = a.pose.position.x - b.pose.position.x;
            let dy = a.pose.position.y - b.pose.position.y;
            let dz = a.pose.position.z - b.pose.position.z;
            (dx * dx + dy * dy + dz * dz).sqrt()
        }

        let mut open_set = HashSet::new();
        let mut came_from = HashMap::new();
        let mut g_score = HashMap::new();
        let mut f_score = HashMap::new();

        open_set.insert(start);
        g_score.insert(start, 0.0);
        f_score.insert(start, heuristic(&start, &goal));

        while !open_set.is_empty() {
            let current = open_set.iter().min_by_key(|&p| f_score[p]).unwrap();
            if *current == goal {
                let mut path = Vec::new();
                let mut current = goal;
                while let Some(prev) = came_from.get(current) {
                    path.push(*prev);
                    current = prev;
                }
                path.reverse();
                return Ok(Path { poses: path });
            }

            open_set.remove(current);

            for neighbor in self.get_neighbors(current) {
                let tentative_g_score = g_score[current] + heuristic(current, &neighbor);
            }
        }

        Err(PathPlannerError::NoPathFound)
    }
}
